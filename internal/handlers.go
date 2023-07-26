package internal

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	cookie, err := r.Cookie("logged-in")

	if err == http.ErrNoCookie || cookie.Value == "not-logged" {
		cookie = &http.Cookie{
			Name:  "logged-in",
			Value: "not-logged",
		}
		http.SetCookie(w, cookie)

		files := []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			return
		}
		tmpl.Execute(w, ShowPost())
	} else if err != nil {
		log.Fatal(err)
	} else {
		if time.Now().After(cookie.Expires) {
			db, err := sql.Open("sqlite3", "./sql/database.db")
			if err != nil {
				log.Println(err.Error())
				return
			}
			tx, err := db.Begin()
			if err != nil {
				log.Fatal(err)
			}
			db.Exec("Delete * from cookies where Id = ( ? )", cookie.Value)
			tx.Commit()
			db.Close()
			cookie = &http.Cookie{
				Name:  "logged-in",
				Value: "not-logged",
			}

		}
		c := cookie.Value
		db, err := sql.Open("sqlite3", "./sql/database.db")
		var name string

		Name, err := db.Query("SELECT lame FROM cookies WHERE Id = ( ? )", c)
		if err != nil {
			log.Fatal(err)
		}
		defer Name.Close()
		for Name.Next() {
			Name.Scan(&name)
			fmt.Println(name)
		}

		files := []string{
			"./ui/html/user.home.tmpl",
			"./ui/html/base.layout.tmpl",
		}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			return
		}
		db.Close()
		tmpl.Execute(w, ShowPost())

	}
}

// 	if err != nil {

// 		log.Fatal(err)
// 		return
// 	}
// 	var name string
// 	Name.Scan(&name)

// 	if r.URL.Path != "/" {
// 		http.NotFound(w, r)

// 		return
// 	}
// 	if name == "" {
// 		if r.Method != http.MethodGet {
// 			w.Header().Set("Allow", http.MethodGet)
// 			w.WriteHeader(http.StatusMethodNotAllowed)
// 			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
// 			return
// 		}
// 		files := []string{
// 			"./ui/html/home.page.tmpl",
// 			"./ui/html/base.layout.tmpl",
// 		}
// 		tmpl, err := template.ParseFiles(files...)
// 		if err != nil {
// 			log.Println(err.Error())
// 			return
// 		}
// 		tmpl.Execute(w, nil)
// 	} else {
// 		files := []string{
// 			"./ui/html/user.home.tmpl",
// 			"./ui/html/base.layout.tmpl",
// 		}
// 		tmpl, err := template.ParseFiles(files...)
// 		if err != nil {
// 			log.Println(err.Error())
// 			return
// 		}
// 		fmt.Println(cookie.Value)

// 		tmpl.Execute(w, name)
// 	}
// }

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	tmpl, err := template.ParseFiles("./ui/html/signup.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	tmpl.Execute(w, nil)
}

func SignUpConfirmation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	name := r.FormValue("UserName")

	email := r.FormValue("UserEmail")
	password := r.FormValue("UserPassword")
	rewrittenPassword := r.FormValue("UserRewrittenPassword")

	result, text := ConfirmSignup(name, email, password, rewrittenPassword)
	if result == true {

		pwd, err := bcrypt.GenerateFromPassword([]byte(password), 1)
		if err != nil {
			log.Fatal(err)
		}
		db, err := sql.Open("sqlite3", "./sql/database.db")
		if err != nil {
			log.Fatal(err)
		}
		AddUser(name, email, string(pwd), db)

		http.Redirect(w, r, "/signin", 302)
	} else {
		tmpl, err := template.ParseFiles("./ui/html/signup.html")
		if err != nil {
			log.Println(err.Error())
			return
		}
		tmpl.Execute(w, text)
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	tmpl, err := template.ParseFiles("./ui/html/signin.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	tmpl.Execute(w, nil)
}

func SignInConfirmation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	name := r.FormValue("UserName")
	password := r.FormValue("UserPassword")
	result, text := ConfirmSignin(name, password)
	if result == true {
		u1, err := uuid.NewV4()
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		u2 := uuid.NewV3(u1, name).String()
		db, err := sql.Open("sqlite3", "./sql/database.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		CreateSession(u2, name, db)

		cookie := &http.Cookie{Name: "logged-in", Value: u2, Expires: time.Now().Add(365 * 24 * time.Hour)}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		tmpl, err := template.ParseFiles("./ui/html/signin.html")
		if err != nil {
			log.Println(err.Error())
			return
		}
		tmpl.Execute(w, text)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}

	tmpl, err := template.ParseFiles("./ui/html/create.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	tmpl.Execute(w, nil)
}

func PostConfirmation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	text := r.FormValue("convert")
	cat := r.FormValue("cars")
	cookie, err := r.Cookie("logged-in")
	if err != nil {
		log.Fatal(err)
	}
	CreatePost(cookie.Value, text, cat)
	http.Redirect(w, r, "/", 302)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}
	cookie, err := r.Cookie("logged-in")
	DeleteCookie(cookie.Value, db)
	cookie = &http.Cookie{
		Name:  "logged-in",
		Value: "not-logged",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}

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

func Filter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	r.ParseForm()
	cookie, err := r.Cookie("logged-in")
	// fmt.Println(r.Form["Category"])
	// fmt.Println(r.Form["LikeDislike"])
	likesdislikes := r.Form["LikeDislike"]
	categories := r.Form["Category"]
	var formattedlikes []string

	for i := range likesdislikes {
		formattedlikes = append(formattedlikes, likesdislikes[i]+"s.Postid")
	}
	// all := append(likes, categories...)
	if len(likesdislikes) == 0 && len(categories) == 0 && len(r.Form["YourPosts"]) == 0 {
		http.Redirect(w, r, "/", 301)
	}

	text := "SELECT posts.Post, posts.Namae, posts.Category, posts.Id from posts "

	if len(likesdislikes) == 2 {
		text = text + "INNER JOIN likes on posts.Id=likes.Postid INNER JOIN dislikes on posts.Id=dislikes.Postid"
	} else if len(likesdislikes) == 1 {
		if likesdislikes[0] == "like" {
			text = text + "INNER JOIN likes on posts.Id=likes.Postid"
		} else {
			text = text + "INNER JOIN dislikes on posts.Id=dislikes.Postid"
		}
	}
	if len(likesdislikes) > 0 {
		// text = text + " WHERE posts.Id IN (" + strings.Join(formattedlikes, ", ") + ")"
		text = text + " WHERE "
		for i := range formattedlikes {
			if i == 0 {
				text = text + "posts.Id=" + formattedlikes[i]
			} else {
				text = text + " OR posts.Id=" + formattedlikes[i]
			}
		}
	} else if len(categories) > 0 {
		text = text + " WHERE "
		for i := range categories {
			// text = text + " OR posts.Category=\"" + categories[i] + "\""
			if i == 0 {
				text = text + "posts.Category=\"" + categories[i] + "\""
			} else {
				text = text + " OR posts.Category=\"" + categories[i] + "\""
			}
		}
	}
	if len(categories) > 0 {
		if len(likesdislikes) > 0 {
			text = text + " AND ("
			for i := range categories {
				// text = text + " OR posts.Category=\"" + categories[i] + "\""
				if i == 0 {
					text = text + "posts.Category=\"" + categories[i] + "\""
				} else {
					text = text + " OR posts.Category=\"" + categories[i] + "\""
				}
			}
			text = text + ")"
		} else {
			text = text + " OR "
			for i := range categories {
				// text = text + " OR posts.Category=\"" + categories[i] + "\""
				if i == 0 {
					text = text + "posts.Category=\"" + categories[i] + "\""
				} else {
					text = text + " OR posts.Category=\"" + categories[i] + "\""
				}
			}

		}
	}
	fmt.Println(text)
	fmt.Println(r.Form["YourPosts"])
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(text)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	var t string
	var n string
	var c string
	var i int
	var likes int
	var dislikes int
	var posts []Post
	var ids []int

	for rows.Next() {
		rows.Scan(&t, &n, &c, &i)
		x := false
		for _, el := range ids {
			if el == i {
				x = true
				break
			}
		}
		if x == true {
			continue
		}
		ids = append(ids, i)

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

	db.Close()

	if len(r.Form["YourPosts"]) == 1 && (len(categories) != 0 || len(likesdislikes) != 0) {
		fmt.Println("G")
		db, err := sql.Open("sqlite3", "./sql/database.db")
		if err != nil {
			log.Fatal(err)
		}
		st, err := db.Query("SELECT lame FROM cookies WHERE Id=(?)", cookie.Value)
		var name string

		for st.Next() {
			st.Scan(&name)
		}
		st.Close()
		var res []Post
		for i := range posts {
			if posts[i].Name == name {
				res = append(res, posts[i])
			}
		}

		files := []string{
			"./ui/html/user.home.tmpl",
			"./ui/html/base.layout.tmpl",
		}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			return
		}
		fmt.Println(res)

		tmpl.Execute(w, res)
		return

	} else if len(r.Form["YourPosts"]) == 1 && len(categories) == 0 && len(likesdislikes) == 0 {
		fmt.Println("GG")
		db, err := sql.Open("sqlite3", "./sql/database.db")
		if err != nil {
			log.Fatal(err)
		}
		st, err := db.Query("SELECT lame FROM cookies WHERE Id=(?)", cookie.Value)
		var name string

		for st.Next() {
			st.Scan(&name)
		}

		st.Close()
		db, err = sql.Open("sqlite3", "./sql/database.db")
		if err != nil {
			log.Fatal(err)
		}
		var res1 []Post
		st1, err := db.Query("SELECT Post,Namae,Category,Id FROM posts WHERE Namae=(?)", name)
		if err != nil {
			log.Fatal(err)
		}
		var t string
		var n string
		var c string
		var i int
		var likes int
		var dislikes int
		defer st1.Close()
		for st1.Next() {
			st1.Scan(&t, &n, &c, &i)

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
			res1 = append(res1, onepost)

		}
		db.Close()

		files := []string{
			"./ui/html/user.home.tmpl",
			"./ui/html/base.layout.tmpl",
		}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			return
		}
		fmt.Print(res1)
		tmpl.Execute(w, res1)
		return
	} else {
		fmt.Println("GGG")

		files := []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
		}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			return
		}

		tmpl.Execute(w, posts)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
func ComDislikes(w http.ResponseWriter, r *http.Request) {}
