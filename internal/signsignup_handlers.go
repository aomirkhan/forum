package internal

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

		pwd, err := bcrypt.GenerateFromPassword([]byte(password), 1)
		if err != nil {
			log.Fatal(err)
		}
		db, err := sql.Open("sqlite3", "./sql/database.db")
		if err != nil {
			log.Fatal(err)
		}
		AddUser(name, email, string(pwd), db)

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
		u1, err := uuid.NewV4()
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		u2 := uuid.NewV3(u1, name).String()
		db, err := sql.Open("sqlite3", "./sql/database.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		CreateSession(u2, name, db)

		cookie := &http.Cookie{Name: "logged-in", Value: u2, Expires: time.Now().Add(365 * 24 * time.Hour)}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		tmpl, err := template.ParseFiles("./ui/html/signin.html")
		if err != nil {
			log.Println(err.Error())
			return
		}
		tmpl.Execute(w, text)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}
	cookie, err := r.Cookie("logged-in")
	DeleteCookie(cookie.Value, db)
	cookie = &http.Cookie{
		Name:  "logged-in",
		Value: "not-logged",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}
