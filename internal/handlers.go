package internal

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	cookie, err := r.Cookie("logged-in")
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "logged-in",
			Value: "not-logged",
		}
		http.SetCookie(w, cookie)

		files := []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
		}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			return
		}
		tmpl.Execute(w, nil)
	} else if err != nil {
		log.Fatal(err)
	} else {
		c := cookie.Value
		fmt.Print(c)
		db, err := sql.Open("sqlite3", "./sql/database.db")
		var name string

		Name, err := db.Query(fmt.Sprintf("SELECT Name FROM cookies WHERE Id =%s", c))
		if err != nil {
			log.Fatal(err)
		}
		Name.Scan(&name)
		files := []string{
			"./ui/html/user.home.tmpl",
			"./ui/html/base.layout.tmpl",
		}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			return
		}
		fmt.Println(cookie.Value)

		tmpl.Execute(w, name)
	}
}

// 	if err != nil {

// 		log.Fatal(err)
// 		return
// 	}
// 	var name string
// 	Name.Scan(&name)

// 	if r.URL.Path != "/" {
// 		http.NotFound(w, r)

// 		return
// 	}
// 	if name == "" {
// 		if r.Method != http.MethodGet {
// 			w.Header().Set("Allow", http.MethodGet)
// 			w.WriteHeader(http.StatusMethodNotAllowed)
// 			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
// 			return
// 		}
// 		files := []string{
// 			"./ui/html/home.page.tmpl",
// 			"./ui/html/base.layout.tmpl",
// 		}
// 		tmpl, err := template.ParseFiles(files...)
// 		if err != nil {
// 			log.Println(err.Error())
// 			return
// 		}
// 		tmpl.Execute(w, nil)
// 	} else {
// 		files := []string{
// 			"./ui/html/user.home.tmpl",
// 			"./ui/html/base.layout.tmpl",
// 		}
// 		tmpl, err := template.ParseFiles(files...)
// 		if err != nil {
// 			log.Println(err.Error())
// 			return
// 		}
// 		fmt.Println(cookie.Value)

// 		tmpl.Execute(w, name)
// 	}
// }

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

		cost, err := bcrypt.Cost([]byte(password))
		if err != nil {
			log.Fatal(err)
		}
		pwd, err := bcrypt.GenerateFromPassword([]byte(password), cost)
		if err != nil {
			log.Fatal(err)
		}
		AddUser(name, email, string(pwd))

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

		CreateSession(u2, name)

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
	// text := r.FormValue("convert")
	// cat := r.FormValue("cars")

	// CreatePost(Name1, text, cat)
	http.Redirect(w, r, "/", 302)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:  "logged-in",
		Value: "0",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}
