package internal

import (
	"database/sql"
	"fmt"
	"log"
)

func AddUser(UserName string, Email string, hashedPassword string, database *sql.DB) {
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (Name VARCHAN(30), Email VARCHAN(45), Password VARCHAN(45))")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	statement, err = db.Prepare("INSERT INTO users (Name, Email,Password) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(UserName, Email, hashedPassword)
	statement.Exec(UserName, Email, hashedPassword)

	defer db.Close()
}

func CreatePost(name string, text string, category string) {
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS posts (Name VARCHAN(30) PRIMARY KEY, PostText VARCHAN(2000), Category VARCHAN(45))")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	statement, err = db.Prepare("INSERT INTO posts (Name, PostText,Category) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(name, text, category)

	defer db.Close()
}
