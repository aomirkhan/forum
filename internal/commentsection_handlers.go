package internal

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func PostPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	xurl := strings.Split(r.URL.String(), "id=")
	if len(xurl) < 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(xurl[1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	if r.URL.String() != "/comments?id="+strconv.Itoa(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	db, err := sql.Open("sqlite3", "./sql/database.db")
	defer db.Close()
	count, err := db.Query("select count(*) from posts;")
	if err != nil {
		log.Fatal(err)
	}
	var i int
	defer count.Close()
	for count.Next() {
		count.Scan(&i)
	}
	if id > i || id < 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tmp, err := template.ParseFiles("./ui/html/comments.html")
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	Comments := CollectComments(id, db)
	err = tmp.Execute(w, Comments)
}

func CommentConfirmation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	text := r.FormValue("comment")
	previousURL := r.Header.Get("Referer")

	xurl := strings.Split(previousURL, "id=")
	if len(xurl) < 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(xurl[1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}
	cookie, err := r.Cookie("logged-in")
	st, err := db.Query("SELECT lame FROM cookies WHERE Id=(?)", cookie.Value)
	if err != nil {
		log.Fatal(err)
	}
	var name string
	for st.Next() {
		st.Scan(&name)
	}
	st.Close()

	AddComment(name, text, id, db)
	http.Redirect(w, r, previousURL, 302)
}
