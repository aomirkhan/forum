package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	gg    sync.Mutex
	name1 string
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}
	myValue := r.Context().Value("myKey").(string)
	fmt.Println(string(myValue))

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	tmpl, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	tmpl.Execute(w, nil)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	tmpl, err := template.ParseFiles("./ui/html/signup.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	tmpl.Execute(w, nil)
}

func SignUpConfirmation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	name := r.FormValue("UserName")

	email := r.FormValue("UserEmail")
	password := r.FormValue("UserPassword")
	rewrittenPassword := r.FormValue("UserRewrittenPassword")

	result, text := ConfirmSignup(name, email, password, rewrittenPassword)
	if result == true {
		db, err := sql.Open("sqlite3", "./sql/database.db")
		if err != nil {
			log.Fatal(err)
		}
		cost, err := bcrypt.Cost([]byte(password))
		pwd, err := bcrypt.GenerateFromPassword([]byte(password), cost)
		AddUser(name, email, string(pwd), db)
		defer db.Close()
		http.Redirect(w, r, "/signin", 302)
	} else {
		tmpl, err := template.ParseFiles("./ui/html/signup.html")
		if err != nil {
			log.Println(err.Error())
			return
		}
		tmpl.Execute(w, text)
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	tmpl, err := template.ParseFiles("./ui/html/signin.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	tmpl.Execute(w, nil)
}

func SignInConfirmation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	name := r.FormValue("UserName")
	password := r.FormValue("UserPassword")
	result, text := ConfirmSignin(name, password)
	if result == true {
		ctx := context.WithValue(r.Context(), "myKey", "aza")
		http.Redirect(w, r.WithContext(ctx), "/", 302)
	} else {
		tmpl, err := template.ParseFiles("./ui/html/signin.html")
		if err != nil {
			log.Println(err.Error())
			return
		}
		tmpl.Execute(w, text)
	}
}

func Feed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	tmpl, err := template.ParseFiles("./ui/html/signin.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	tmpl.Execute(w, nil)
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}

	tmpl, err := template.ParseFiles("./ui/html/create.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	tmpl.Execute(w, nil)
}

func PostConfirmation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	text := r.FormValue("convert")
	cat := r.FormValue("cars")

	CreatePost(name1, text, cat)
	http.Redirect(w, r, "/", 302)
}
