package main

import (
	"log"
	"net/http"
	"text/template"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.Write([]byte("Rizalox"))
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
	}
	tmpl, err := template.ParseFiles("./templates/home.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	tmpl.Execute(w, nil)
}
