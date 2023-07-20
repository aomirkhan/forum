package internal

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func AddUser(UserName string, Email string, hashedPassword string, db *sql.DB) {
	statement, err := db.Prepare("INSERT INTO users (Name, Email,Password) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(UserName, Email, hashedPassword)
	statement.Exec(UserName, Email, hashedPassword)
	db.Close()
}

func CreatePost(cookie string, text string, category string) {
	db, err := sql.Open("sqlite3", "./sql/database.db")
	Name, err := db.Query("SELECT lame FROM cookies WHERE Id = ( ? )", cookie)
	if err != nil {
		log.Fatal(err)
	}
	defer Name.Close()
	var name string
	for Name.Next() {
		Name.Scan(&name)
	}
	if err != nil {
		log.Fatal(err)
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	Flag, err := db.Query("SELECT count(*) FROM posts")
	defer Flag.Close()
	var flag int
	for Flag.Next() {
		Flag.Scan(&flag)
	}

	_, err = db.Exec("INSERT INTO posts (Post,Namae,Category,Id) VALUES (?, ?, ?, ? )", text, name, category, flag+1)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	db.Close()
}

func CreateSession(id, name string, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO cookies (Id, lame) VALUES (?, ?)", id, name)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	db.Close()
}

func DeleteCookie(cookie string, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM cookies WHERE Id=(?)", cookie)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	db.Close()
}

func AddComment(name, text string, id int, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO comments (Name,Text,Id) VALUES (?, ?, ?)", name, text, id)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	db.Close()
}
