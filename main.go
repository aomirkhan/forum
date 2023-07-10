package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Homepage)
	fmt.Println("http://127.0.0.1:8000")
	http.ListenAndServe(":8000", mux)
}
