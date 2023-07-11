package internal

import (
	"database/sql"
	"log"
)

func AddUser(UserName string, Email string, hashedPassword string, database *sql.DB) {
	db, err := sql.Open("sqlite3", "./example.db")
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
	statement.Exec(UserName, Email, hashedPassword)

	defer db.Close()
}
