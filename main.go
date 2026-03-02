package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

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
	dst, err := os.Create("./audioFiles/" + time.Now().String() + fileHeader.Filename)
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
	cmd := exec.Command("python", "LoadAudio.py")
	if err := cmd.Start(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := cmd.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("cargando audio..."))
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
	response, err := http.Get("http://127.0.0.1:8090?prompt=" + prompt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	w.Write(body)
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/audio", handleFile)
	http.HandleFunc("/prompt", handlePrompt)
	http.HandleFunc("/", handleHome)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
