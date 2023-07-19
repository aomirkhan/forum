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

func CreatePost(name string, text string, category string, db *sql.DB) {
	statement, err := db.Prepare("INSERT INTO posts (Name, PostText,Category) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(name, text, category)

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
