package main

import (
	"fmt"
	"forum/internal"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("./templates"))
	mux.Handle("/templates/", http.StripPrefix("/templates", files))

	mux.HandleFunc("/", internal.Homepage)
	mux.HandleFunc("/signup", internal.SignUp)
	mux.HandleFunc("/signin", internal.SignIn)
	mux.HandleFunc("/signupconfirmation", internal.SignUpConfirmation)
	mux.HandleFunc("/signinconfirmation", internal.SignInConfirmation)
	fmt.Println("http://127.0.0.1:8000")
	http.ListenAndServe(":8000", mux)
}
