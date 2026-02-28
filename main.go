package main

import (
	"html/template"
	"log"
	"net/http"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("htmls/main.html"))
	if err := t.Execute(w, nil); err != nil {
		log.Fatal(err)
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
