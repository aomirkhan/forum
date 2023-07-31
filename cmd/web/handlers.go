package web

import (
	"database/sql"
	"fmt"
	"forum/internal"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusNotFound)
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
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, internal.ShowPost())
	} else if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	} else {
		if time.Now().After(cookie.Expires) {
			db, err := sql.Open("sqlite3", "./sql/database.db")
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			tx, err := db.Begin()
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
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
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		defer Name.Close()
		for Name.Next() {
			Name.Scan(&name)
		}

		files := []string{
			"./ui/html/user.home.tmpl",
			"./ui/html/base.layout.tmpl",
		}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		db.Close()
		tmpl.Execute(w, internal.ShowPost())

	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed)
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
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	name := r.FormValue("UserName")

	email := r.FormValue("UserEmail")
	password := r.FormValue("UserPassword")
	rewrittenPassword := r.FormValue("UserRewrittenPassword")

	result, text := internal.ConfirmSignup(name, email, password, rewrittenPassword)
	if result == true {

		pwd, err := bcrypt.GenerateFromPassword([]byte(password), 1)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		db, err := sql.Open("sqlite3", "./sql/database.db")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		internal.AddUser(name, email, string(pwd), db)

		http.Redirect(w, r, "/signin", 302)
	} else {
		tmpl, err := template.ParseFiles("./ui/html/signup.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, text)
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	tmpl, err := template.ParseFiles("./ui/html/signin.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func SignInConfirmation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	name := r.FormValue("UserName")
	password := r.FormValue("UserPassword")
	result, text := internal.ConfirmSignin(name, password)
	if result == true {
		u1, err := uuid.NewV4()
		if err != nil {
			ErrorHandler(w, http.StatusForbidden)
			return
		}
		u2 := uuid.NewV3(u1, name).String()
		db, err := sql.Open("sqlite3", "./sql/database.db")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		defer db.Close()
		internal.CreateSession(u2, name, db)

		cookie := &http.Cookie{Name: "logged-in", Value: u2, Expires: time.Now().Add(365 * 24 * time.Hour)}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		tmpl, err := template.ParseFiles("./ui/html/signin.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, text)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	cookie, err := r.Cookie("logged-in")
	internal.DeleteCookie(cookie.Value, db)
	cookie = &http.Cookie{
		Name:  "logged-in",
		Value: "not-logged",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed)
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
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	title := r.FormValue("title")
	text := r.FormValue("convert")
	cat := r.FormValue("cars")
	cookie, err := r.Cookie("logged-in")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	internal.CreatePost(cookie.Value, text, cat, title)
	http.Redirect(w, r, "/", 302)
}

func PostPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	xurl := strings.Split(r.URL.String(), "id=")
	if len(xurl) < 2 {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(xurl[1])
	if err != nil {
		ErrorHandler(w, http.StatusNotFound)
	}
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	qu, err := db.Query("select Title, Post,Namae from posts where Id=(?)", id)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	defer qu.Close()
	var title string
	var text string
	var name string
	for qu.Next() {
		qu.Scan(&title, &text, &name)
	}

	db.Close()
	if r.URL.String() != "/comments?id="+strconv.Itoa(id) {
		ErrorHandler(w, http.StatusNotFound)
		return
	}
	db, err = sql.Open("sqlite3", "./sql/database.db")
	defer db.Close()
	count, err := db.Query("select count(*) from posts;")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	var i int
	defer count.Close()
	for count.Next() {
		count.Scan(&i)
	}
	if id > i || id < 1 {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	tmp, err := template.ParseFiles("./ui/html/comments.html")
	if err != nil {

		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	comments := internal.CollectComments(id, db)
	result := internal.Postpage{
		Title:    title,
		Post:     text,
		Name:     name,
		Comments: comments,
	}
	// fmt.Printf("%s i title\n%s is post\n%s is name\n", title, text, name)
	// fmt.Println(comments)
	err = tmp.Execute(w, result)
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
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	cookie, err := r.Cookie("logged-in")
	st, err := db.Query("SELECT lame FROM cookies WHERE Id=(?)", cookie.Value)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	var name string
	for st.Next() {
		st.Scan(&name)
	}
	st.Close()

	internal.AddComment(name, text, id, db)
	http.Redirect(w, r, previousURL, 302)
}

func Filter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {

		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	r.ParseForm()
	cookie, err := r.Cookie("logged-in")
	cc := cookie.Value
	db, err := sql.Open("sqlite3", "./sql/database.db")
	var namecookie string

	Name, err := db.Query("SELECT lame FROM cookies WHERE Id = ( ? )", cc)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	defer Name.Close()
	for Name.Next() {
		Name.Scan(&namecookie)
	}
	Name.Close()
	db.Close()

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

	text := "SELECT posts.Title, posts.Post, posts.Namae, posts.Category, posts.Id from posts "

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
				text = text + "(posts.Id=" + formattedlikes[i] + " AND " + likesdislikes[i] + "s.Name=\"" + namecookie + "\")"
			} else {
				text = text + " OR (posts.Id=" + formattedlikes[i] + " AND " + likesdislikes[i] + "s.Name=\"" + namecookie + "\")"
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

	db, err = sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(text)
	defer db.Close()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	var title string
	var t string
	var n string
	var c string
	var i int
	var likes int
	var dislikes int
	var posts []internal.Post
	var ids []int

	for rows.Next() {
		rows.Scan(&title, &t, &n, &c, &i)
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
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		err = db.QueryRow("SELECT count(*) FROM dislikes WHERE Postid=(?)", i).Scan(&dislikes)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		onepost := internal.Post{
			Title:    title,
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
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		st, err := db.Query("SELECT lame FROM cookies WHERE Id=(?)", cookie.Value)
		var name string

		for st.Next() {
			st.Scan(&name)
		}
		st.Close()
		var res []internal.Post
		for i := range posts {
			if posts[i].Name == name {
				res = append(res, posts[i])
			}
		}
		cook, err := r.Cookie("logged-in")
		var files []string
		if err == http.ErrNoCookie || cook.Value == "not-logged" {
			files = []string{
				"./ui/html/home.page.tmpl",
				"./ui/html/base.layout.tmpl",
			}
		} else {
			files = []string{
				"./ui/html/user.home.tmpl",
				"./ui/html/base.layout.tmpl",
			}
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, res)
		return

	} else if len(r.Form["YourPosts"]) == 1 && len(categories) == 0 && len(likesdislikes) == 0 {
		fmt.Println("GG")
		db, err := sql.Open("sqlite3", "./sql/database.db")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		st, err := db.Query("SELECT lame FROM cookies WHERE Id=(?)", cookie.Value)
		var name string

		for st.Next() {
			st.Scan(&name)
		}

		st.Close()
		db, err = sql.Open("sqlite3", "./sql/database.db")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		var res1 []internal.Post
		st1, err := db.Query("SELECT Title, Post,Namae,Category,Id FROM posts WHERE Namae=(?)", name)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		var title string
		var t string
		var n string
		var c string
		var i int
		var likes int
		var dislikes int
		defer st1.Close()
		for st1.Next() {
			st1.Scan(&title, &t, &n, &c, &i)

			err := db.QueryRow("SELECT count(*) FROM likes WHERE Postid=(?)", i).Scan(&likes)
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			err = db.QueryRow("SELECT count(*) FROM dislikes WHERE Postid=(?)", i).Scan(&dislikes)
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			onepost := internal.Post{
				Title:    title,
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

		cook, err := r.Cookie("logged-in")
		var files []string
		if err == http.ErrNoCookie || cook.Value == "not-logged" {
			files = []string{
				"./ui/html/home.page.tmpl",
				"./ui/html/base.layout.tmpl",
			}
		} else {
			files = []string{
				"./ui/html/user.home.tmpl",
				"./ui/html/base.layout.tmpl",
			}
		}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, res1)

	} else {

		cook, err := r.Cookie("logged-in")
		var files []string
		if err == http.ErrNoCookie || cook.Value == "not-logged" {
			files = []string{
				"./ui/html/home.page.tmpl",
				"./ui/html/base.layout.tmpl",
			}
		} else {
			files = []string{
				"./ui/html/user.home.tmpl",
				"./ui/html/base.layout.tmpl",
			}
		}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, posts)
	}
}

func Likes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	id := r.FormValue("id")

	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	cookie, err := r.Cookie("logged-in")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	var checkName string
	row, err := db.Query("SELECT lame FROM cookies WHERE Id=(?)", cookie.Value)
	for row.Next() {
		row.Scan(&checkName)
	}
	row.Close()
	db, err = sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
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
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO likes (Name, Postid) VALUES (?, ?)", checkName, id)
		_, err = db.Exec("DELETE FROM dislikes WHERE Name=(?) and Postid=(?)", checkName, id)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tx.Commit()
		db.Close()
	} else if checklikes == false && checkdislikes == false {
		tx, err := db.Begin()
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO likes (Name, Postid) VALUES (?, ?)", checkName, id)

		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tx.Commit()
		db.Close()
	} else if checklikes == true && checkdislikes == false {
		tx, err := db.Begin()
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("DELETE FROM likes WHERE Name=(?) and Postid=(?)", checkName, id)

		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tx.Commit()
		db.Close()
	}

	http.Redirect(w, r, previousURL, 302)
}

func Dislikes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	id := r.FormValue("id")

	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	cookie, err := r.Cookie("logged-in")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	var checkName string
	row, err := db.Query("SELECT lame FROM cookies WHERE Id=(?)", cookie.Value)
	for row.Next() {
		row.Scan(&checkName)
	}
	row.Close()
	db, err = sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
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
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO dislikes (Name, Postid) VALUES (?, ?)", checkName, id)
		_, err = db.Exec("DELETE FROM likes WHERE Name=(?) and Postid = (?)", checkName, id)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tx.Commit()
		db.Close()
	} else if checklikes == false && checkdislikes == false {
		tx, err := db.Begin()
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO dislikes (Name, Postid) VALUES (?, ?)", checkName, id)

		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tx.Commit()
		db.Close()
	} else if checklikes == false && checkdislikes == true {
		tx, err := db.Begin()
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("DELETE FROM dislikes WHERE Name=(?) and Postid=(?)", checkName, id)

		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tx.Commit()
		db.Close()

	}

	http.Redirect(w, r, previousURL, 302)
}

func ComLikes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	previousURL := r.Header.Get("Referer")
	postid := (strings.Split(previousURL, "id="))[1]
	id := r.FormValue("id")

	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	cookie, err := r.Cookie("logged-in")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	var checkName string
	row, err := db.Query("SELECT lame FROM cookies WHERE Id=(?)", cookie.Value)
	for row.Next() {
		row.Scan(&checkName)
	}
	row.Close()
	db, err = sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	checklikes := false
	checkdislikes := false
	rows, err := db.Query("SELECT Name FROM comlikes WHERE (Comid,Id)=(?,?)", id, postid)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
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
		ErrorHandler(w, http.StatusInternalServerError)
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
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	fmt.Println("GG", id)
	if checklikes == false && checkdislikes == true {

		_, err = db.Exec("INSERT INTO comlikes (Name, Comid,Id) VALUES (?, ?, ?)", checkName, id, postid)
		_, err = db.Exec("DELETE FROM comdislikes WHERE Name=(?) and Comid=(?) and Id=(?)", checkName, id, postid)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

	} else if checklikes == false && checkdislikes == false {

		_, err = db.Exec("INSERT INTO comlikes (Name, Comid,Id) VALUES (?, ?, ?)", checkName, id, postid)

		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

	} else if checklikes == true && checkdislikes == false {

		_, err = db.Exec("DELETE FROM comlikes WHERE Name=(?) and Comid=(?) and Id=(?)", checkName, id, postid)

		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()

	http.Redirect(w, r, previousURL, 302)
}

func ComDislikes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	previousURL := r.Header.Get("Referer")
	postid := (strings.Split(previousURL, "id="))[1]
	id := r.FormValue("id")

	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	cookie, err := r.Cookie("logged-in")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	var checkName string
	row, err := db.Query("SELECT lame FROM cookies WHERE Id=(?)", cookie.Value)
	for row.Next() {
		row.Scan(&checkName)
	}
	row.Close()
	db, err = sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	checklikes := false
	checkdislikes := false
	rows, err := db.Query("SELECT Name FROM comlikes WHERE (Comid,Id)=(?,?)", id, postid)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
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
		ErrorHandler(w, http.StatusInternalServerError)
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
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	if checklikes == true && checkdislikes == false {

		_, err = db.Exec("INSERT INTO comdislikes (Name, Comid,Id) VALUES (?, ?, ?)", checkName, id, postid)
		_, err = db.Exec("DELETE FROM comlikes WHERE Name=(?) and Comid=(?) and Id=(?)", checkName, id, postid)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

	} else if checklikes == false && checkdislikes == false {

		_, err = db.Exec("INSERT INTO comdislikes (Name, Comid,Id) VALUES (?, ?, ?)", checkName, id, postid)

		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

	} else if checklikes == false && checkdislikes == true {

		_, err = db.Exec("DELETE FROM comdislikes WHERE Name=(?) and Comid=(?) and Id=(?)", checkName, id, postid)

		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()

	http.Redirect(w, r, previousURL, 302)
}

func ErrorHandler(w http.ResponseWriter, status int) {
	tmp, err := template.ParseFiles("./ui/html/error.html")
	if err != nil || tmp == nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	var Err internal.ErrorStruct
	Err.Message = http.StatusText(status)
	Err.Status = status
	err = tmp.Execute(w, Err)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
}
