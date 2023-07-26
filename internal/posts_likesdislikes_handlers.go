package internal

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func Likes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	id := r.FormValue("id")
	fmt.Println("GG", id)
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}
	cookie, err := r.Cookie("logged-in")
	if err != nil {
		log.Fatal(err)
	}
	var checkName string
	row, err := db.Query("SELECT lame FROM cookies WHERE Id=(?)", cookie.Value)
	for row.Next() {
		row.Scan(&checkName)
	}
	row.Close()
	db, err = sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}

	previousURL := r.Header.Get("Referer")
	checklikes := false
	checkdislikes := false
	rows, err := db.Query("SELECT Name FROM likes WHERE Postid=(?)", id)
	var likerName string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&likerName)
		if likerName == checkName {
			checklikes = true
		}
	}
	x, err := db.Query("SELECT Name FROM dislikes WHERE Postid=(?)", id)
	var dislikerName string
	defer x.Close()

	for x.Next() {
		x.Scan(&dislikerName)
		if dislikerName == checkName {
			checkdislikes = true
		}
	}
	if checklikes == false && checkdislikes == true {

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("INSERT INTO likes (Name, Postid) VALUES (?, ?)", checkName, id)
		_, err = db.Exec("DELETE FROM dislikes WHERE Name=(?) and Postid=(?)", checkName, id)
		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
		db.Close()
	} else if checklikes == false && checkdislikes == false {
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("INSERT INTO likes (Name, Postid) VALUES (?, ?)", checkName, id)

		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
		db.Close()
	} else if checklikes == true && checkdislikes == false {
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("DELETE FROM likes WHERE Name=(?) and Postid=(?)", checkName, id)

		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
		db.Close()
	}

	http.Redirect(w, r, previousURL, 302)
}

func Dislikes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	id := r.FormValue("id")
	fmt.Println("GG", id)
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}
	cookie, err := r.Cookie("logged-in")
	if err != nil {
		log.Fatal(err)
	}
	var checkName string
	row, err := db.Query("SELECT lame FROM cookies WHERE Id=(?)", cookie.Value)
	for row.Next() {
		row.Scan(&checkName)
	}
	row.Close()
	db, err = sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}

	previousURL := r.Header.Get("Referer")
	checklikes := false
	checkdislikes := false
	rows, err := db.Query("SELECT Name FROM likes WHERE Postid=(?)", id)
	var likerName string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&likerName)
		if likerName == checkName {
			checklikes = true
		}
	}
	x, err := db.Query("SELECT Name FROM dislikes WHERE Postid=(?)", id)
	var dislikerName string
	defer x.Close()

	for x.Next() {
		x.Scan(&dislikerName)
		if dislikerName == checkName {
			checkdislikes = true
		}
	}
	if checklikes == true && checkdislikes == false {

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("INSERT INTO dislikes (Name, Postid) VALUES (?, ?)", checkName, id)
		_, err = db.Exec("DELETE FROM likes WHERE Name=(?) and Postid = (?)", checkName, id)
		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
		db.Close()
	} else if checklikes == false && checkdislikes == false {
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("INSERT INTO dislikes (Name, Postid) VALUES (?, ?)", checkName, id)

		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
		db.Close()
	} else if checklikes == false && checkdislikes == true {
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("DELETE FROM dislikes WHERE Name=(?) and Postid=(?)", checkName, id)

		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
		db.Close()

	}

	http.Redirect(w, r, previousURL, 302)
}
