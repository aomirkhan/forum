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
	Likes    int
	Dislikes int
	// Comments [string]string
}

type Comment struct {
	Name  string
	Text  string
	ComId int
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
	var likes int
	var dislikes int
	defer row.Close()
	for row.Next() {
		row.Scan(&t, &n, &c, &i)
		err := db.QueryRow("SELECT count(*) FROM likes WHERE Postid=(?)", i).Scan(&likes)
		if err != nil {
			log.Fatal(err)
		}
		err = db.QueryRow("SELECT count(*) FROM dislikes WHERE Postid=(?)", i).Scan(&dislikes)
		if err != nil {
			log.Fatal(err)
		}
		onepost := Post{
			Text:     t,
			Name:     n,
			Category: c,
			Id:       i,
			Likes:    likes,
			Dislikes: dislikes,
		}
		posts = append(posts, onepost)
	}

	tx.Commit()
	db.Close()
	return posts
}
