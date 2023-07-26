package main

import (
	"fmt"
	"forum/cmd/web"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("./templates"))
	mux.Handle("/templates/", http.StripPrefix("/templates", files))

	mux.HandleFunc("/", web.Homepage)
	mux.HandleFunc("/signup", web.SignUp)
	mux.HandleFunc("/signin", web.SignIn)
	mux.HandleFunc("/logout", web.Logout)
	mux.HandleFunc("/signupconfirmation", web.SignUpConfirmation)
	mux.HandleFunc("/signinconfirmation", web.SignInConfirmation)
	mux.HandleFunc("/comments", web.PostPage)
	mux.HandleFunc("/postconfirmation", web.PostConfirmation)
	mux.HandleFunc("/commentconfirmation", web.CommentConfirmation)
	mux.HandleFunc("/create", web.Create)
	mux.HandleFunc("/like", web.Likes)
	mux.HandleFunc("/dislike", web.Dislikes)
	mux.HandleFunc("/filter", web.Filter)
	mux.HandleFunc("/comlike", web.ComLikes)
	mux.HandleFunc("/comdislike", web.ComDislikes)
	mux.HandleFunc("/filter/likes", web.Likes)
	fmt.Println("http://127.0.0.1:8000")
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}
