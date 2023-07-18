package internal

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func AddUser(UserName string, Email string, hashedPassword string) {
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}
	// statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (Name VARCHAN(30), Email VARCHAN(45), Password VARCHAN(45))")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// statement.Exec()
	statement, err := db.Prepare("INSERT INTO users (Name, Email,Password) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(UserName, Email, hashedPassword)
	statement.Exec(UserName, Email, hashedPassword)
	db.Close()
}

func CreatePost(name string, text string, category string) {
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}
	// statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS posts (Name VARCHAN(30) PRIMARY KEY, PostText VARCHAN(2000), Category VARCHAN(45))")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// statement.Exec()
	statement, err := db.Prepare("INSERT INTO posts (Name, PostText,Category) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(name, text, category)
	db.Close()
}

func CreateSession(id, name string) {
	fmt.Println(id)
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}
	// statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS cookies (Id VARCHAR(100) PRIMARY KEY, Name VARCHAR(30))")
	// if err != nil {
	// 	fmt.Println("Rizachert")
	// 	log.Fatal(err)
	// }
	// statement.Exec()
	statement, err := db.Prepare("INSERT INTO cookies (Name,Id) VALUES (?,?)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(name, id)
	rows, err := db.Query("SELECT Name,Id from cookies")
	for rows.Next() {
		var n string
		var i string
		rows.Scan(&n, &i)
		fmt.Printf("Name is %s, ID is %s\n", n, i)
	}
	db.Close()
}
