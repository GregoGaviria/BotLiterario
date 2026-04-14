package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var doneflag = false
var processingflag = false

func handleHome(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("htmls/main.html"))
	if err := t.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}

func handleFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()
	filename := "./audioFiles/" + time.Now().String() + fileHeader.Filename
	dst, err := os.Create(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	go func(filename string) {
		if processingflag == true {
			log.Println("already loading")
			return
		}
		doneflag = false
		processingflag = true

		cmd := exec.Command("python", "LoadAudio.py", "-f", filename)
		// cmd := exec.Command("whisper", filename, "--model", "turbo")

		// cmd.Stdout = buff
		// cmd.Stderr = buff

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// output, err := cmd.CombinedOutput()
		// if err != nil {
		// 	http.Error(w, err.Error()+" "+string(output), http.StatusInternalServerError)
		// 	return
		// }
		// log.Println(string(output))

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("completo")
		doneflag = true
		processingflag = false
	}(filename)

	// w.Write([]byte("cargando audio..."))
	t := template.Must(template.ParseFiles("htmls/loading.html"))
	dots = 1
	if err := t.Execute(w, "⋯"); err != nil {
		log.Fatal(err)
	}
}

var dots int

func handleProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !doneflag {
		var icon string
		switch dots {
		case 0:
			icon = "⋯"
			dots++
		case 1:
			icon = "⋱"
			dots++
		case 2:
			icon = "⋮"
			dots++
		case 3:
			icon = "⋰"
			dots = 0
		}
		t := template.Must(template.ParseFiles("htmls/loading.html"))
		if err := t.Execute(w, icon); err != nil {
			log.Fatal(err)
		}
	} else {
		file, err := os.ReadFile("transcript.txt")
		if err != nil {
			log.Fatal(err)
		}
		reqbody := struct {
			Transcript string `json:"transcript"`
		}{Transcript: string(file)}
		jsonByte, err := json.Marshal(reqbody)
		if err != nil {
			log.Fatal(err)
		}
		response, err := http.Post(
			"http://127.0.0.1:5000/set-context",
			"application/json",
			bytes.NewBuffer(jsonByte),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bodyStruct := struct {
			Characters int    `json:"characters"`
			Message    string `json:"message"`
			Ok         bool   `json:"ok"`
			Session_id string `json:"session_id"`
		}{}
		if err := json.Unmarshal(body, &bodyStruct); err != nil {
			log.Fatal(err)
		}
		response.Body.Close()
		if response.StatusCode > 299 {
			errstring := fmt.Sprintf(
				"statuscode %d: body: %s",
				response.StatusCode,
				string(body),
			)
			http.Error(w, errstring, http.StatusInternalServerError)
			return
		}

		sessionId = bodyStruct.Session_id
		log.Println(sessionId)
		log.Println(bodyStruct.Characters)
		log.Println(bodyStruct.Ok)
		log.Println(bodyStruct.Message)

		t := template.Must(template.ParseFiles("htmls/loadprompt.html"))
		if err := t.Execute(w, nil); err != nil {
			log.Fatal(err)
		}
	}
}

var sessionId string

func handleLoadPrompt(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	t := template.Must(template.ParseFiles("htmls/prompt.html"))
	if err := t.Execute(w, nil); err != nil {
		log.Fatal(err)
	}

}

func handlePrompt(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	prompt := r.FormValue("prompt")
	if prompt == "" {
		http.Error(w, "no hay prompt", http.StatusBadRequest)
		return
	}
	reqbody := struct {
		Question   string `json:"question"`
		Session_id string `json:"session_id"`
	}{Question: prompt, Session_id: sessionId}
	log.Println(reqbody.Question)
	log.Println(reqbody.Session_id)
	jsonByte, err := json.Marshal(reqbody)
	if err != nil {
		log.Fatal(err)
	}
	response, err := http.Post(
		"http://127.0.0.1:5000/ask",
		"application/json",
		bytes.NewBuffer(jsonByte),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bodyStruct := struct {
		Answer     string    `json:"answer"`
		Model      string `json:"model"`
		Ok         bool   `json:"ok"`
		Session_id string `json:"session_id"`
	}{}
	if err := json.Unmarshal(body, &bodyStruct); err != nil {
		log.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode > 299 {
		errstring := fmt.Sprintf(
			"statuscode %d: body: %s",
			response.StatusCode,
			string(body),
		)
		http.Error(w, errstring, http.StatusInternalServerError)
		return
	}

	log.Println(bodyStruct.Answer)
	log.Println(bodyStruct.Ok)
	log.Println(bodyStruct.Model)

	t := template.Must(template.ParseFiles("htmls/message.html"))
	if err := t.Execute(
		w,
		struct{
			U string
			C string
		}{
			U: reqbody.Question,
			C:bodyStruct.Answer,
		},
	); err != nil {
		log.Fatal(err)
	}
	// response.Body.Close()
	// if response.StatusCode > 299 {
	// 	errstring := fmt.Sprintf(
	// 		"statuscode %d: body: %s",
	// 		response.StatusCode,
	// 		string(body),
	// 	)
	// 	http.Error(w, errstring, http.StatusInternalServerError)
	// 	return
	// }
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/audio", handleFile)
	http.HandleFunc("/prompt", handlePrompt)
	http.HandleFunc("/progress", handleProgress)
	http.HandleFunc("/loadPrompt", handleLoadPrompt)
	http.HandleFunc("/", handleHome)
	port := ":8000"
	log.Print("corriendo servidor en puerto " + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
