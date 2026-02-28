package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
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
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()
	dst, err := os.Create("./audioFiles/" + time.Now().String()+fileHeader.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_,err = io.Copy(dst,file)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handleHome)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
