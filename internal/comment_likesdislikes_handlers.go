package internal

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func ComLikes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	previousURL := r.Header.Get("Referer")
	postid := (strings.Split(previousURL, "id="))[1]
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
	checklikes := false
	checkdislikes := false
	rows, err := db.Query("SELECT Name FROM comlikes WHERE (Comid,Id)=(?,?)", id, postid)
	if err != nil {
		log.Fatal(err)
	}
	var likerName string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&likerName)
		if likerName == checkName {
			checklikes = true
		}
	}
	x, err := db.Query("SELECT Name FROM comdislikes WHERE (Comid,Id)=(?,?)", id, postid)
	if err != nil {
		fmt.Println("2")
		return
	}

	var dislikerName string
	defer x.Close()

	for x.Next() {
		x.Scan(&dislikerName)
		if dislikerName == checkName {
			checkdislikes = true
		}
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GG", id)
	if checklikes == false && checkdislikes == true {

		_, err = db.Exec("INSERT INTO comlikes (Name, Comid,Id) VALUES (?, ?, ?)", checkName, id, postid)
		_, err = db.Exec("DELETE FROM comdislikes WHERE Name=(?) and Comid=(?) and Id=(?)", checkName, id, postid)
		if err != nil {
			log.Fatal(err)
		}

	} else if checklikes == false && checkdislikes == false {
		fmt.Println(1)

		_, err = db.Exec("INSERT INTO comlikes (Name, Comid,Id) VALUES (?, ?, ?)", checkName, id, postid)

		if err != nil {
			log.Fatal(err)
		}

	} else if checklikes == true && checkdislikes == false {
		fmt.Println(2)

		_, err = db.Exec("DELETE FROM comlikes WHERE Name=(?) and Comid=(?) and Id=(?)", checkName, id, postid)

		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	http.Redirect(w, r, previousURL, 302)
}

func ComDislikes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("WHY")
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	previousURL := r.Header.Get("Referer")
	postid := (strings.Split(previousURL, "id="))[1]
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
	checklikes := false
	checkdislikes := false
	rows, err := db.Query("SELECT Name FROM comlikes WHERE (Comid,Id)=(?,?)", id, postid)
	if err != nil {
		log.Fatal(err)
	}
	var likerName string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&likerName)
		if likerName == checkName {
			checklikes = true
		}
	}
	x, err := db.Query("SELECT Name FROM comdislikes WHERE (Comid,Id)=(?,?)", id, postid)
	if err != nil {
		fmt.Println("2")
		return
	}

	var dislikerName string
	defer x.Close()

	for x.Next() {
		x.Scan(&dislikerName)
		if dislikerName == checkName {
			checkdislikes = true
		}
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GG", id)
	if checklikes == true && checkdislikes == false {

		_, err = db.Exec("INSERT INTO comdislikes (Name, Comid,Id) VALUES (?, ?, ?)", checkName, id, postid)
		_, err = db.Exec("DELETE FROM comlikes WHERE Name=(?) and Comid=(?) and Id=(?)", checkName, id, postid)
		if err != nil {
			log.Fatal(err)
		}

	} else if checklikes == false && checkdislikes == false {
		fmt.Println(1)

		_, err = db.Exec("INSERT INTO comdislikes (Name, Comid,Id) VALUES (?, ?, ?)", checkName, id, postid)

		if err != nil {
			log.Fatal(err)
		}

	} else if checklikes == false && checkdislikes == true {
		fmt.Println(2)

		_, err = db.Exec("DELETE FROM comdislikes WHERE Name=(?) and Comid=(?) and Id=(?)", checkName, id, postid)

		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	http.Redirect(w, r, previousURL, 302)
}
