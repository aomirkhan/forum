package internal

import (
	"database/sql"
	"log"
)

type Post struct {
	Text     string
	Name     string
	Category string
	Id       int
	// Comments [string]string
}

type Comment struct {
	Name string
	Text string
}

func ShowPost() []Post {
	var posts []Post
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	row, err := db.Query("SELECT * FROM posts")
	if err != nil {
		log.Fatal(err)
	}
	var t string
	var n string
	var c string
	var i int
	defer row.Close()
	for row.Next() {
		row.Scan(&t, &n, &c, &i)
		onepost := Post{
			Text:     t,
			Name:     n,
			Category: c,
			Id:       i,
		}
		posts = append(posts, onepost)
	}

	tx.Commit()
	db.Close()
	return posts
}
