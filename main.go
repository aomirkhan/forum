package main

import (
	"fmt"
	"forum/internal"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("./templates"))
	mux.Handle("/templates/", http.StripPrefix("/templates", files))

	mux.HandleFunc("/", internal.Homepage)
	mux.HandleFunc("/signup", internal.SignUp)
	mux.HandleFunc("/signin", internal.SignIn)
	mux.HandleFunc("/logout", internal.Logout)
	mux.HandleFunc("/signupconfirmation", internal.SignUpConfirmation)
	mux.HandleFunc("/signinconfirmation", internal.SignInConfirmation)
	mux.HandleFunc("/comments", internal.PostPage)
	mux.HandleFunc("/postconfirmation", internal.PostConfirmation)
	mux.HandleFunc("/commentconfirmation", internal.CommentConfirmation)
	mux.HandleFunc("/create", internal.Create)
	mux.HandleFunc("/like", internal.Likes)
	mux.HandleFunc("/dislike", internal.Dislikes)
	mux.HandleFunc("/filter", internal.Filter)
	mux.HandleFunc("/comlike", internal.ComLikes)
	mux.HandleFunc("/comdislike", internal.ComDislikes)
	mux.HandleFunc("/filter/likes", internal.Likes)
	fmt.Println("http://127.0.0.1:8000")
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}
