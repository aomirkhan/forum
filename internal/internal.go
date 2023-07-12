package internal

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var nameuser chan string

func ConfirmSignup(Name string, Email string, Password string, RewrittenPassword string) (bool, string) {
	if RewrittenPassword != Password {
		return false, "Passwords don't match, write again."
	}
	if len(Name) < 3 || len(Password) < 7 {
		return false, "Name/Password doesn't have enough characters. Minimum for name is 3 and for password is 7."
	}
	if isEmailValid(Email) == false {
		return false, "Wrong format for Email"
	}
	db, err := sql.Open("sqlite3", "./sql/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT Name FROM users")
	row1 := db.QueryRow(query)

	var name string
	err = row1.Scan(&name)

	if err == sql.ErrNoRows {
	} else if err != nil {
		log.Fatal(err)
	} else {

		rows, err := db.Query("SELECT Name, Email FROM users")
		if err != nil {
			log.Fatal(err)
		}
		var name string
		var email string
		for rows.Next() {
			rows.Scan(&name, &email)
			if name == Name {
				return false, "That name is already being used"
			} else if Email == email {
				return false, "That Email is already being used"
			}
		}
	}

	return true, "OK"
}

func ConfirmSignin(Name string, Password string) (bool, string) {
	db, _ := sql.Open("sqlite3", "./sql/database.db")

	rows, err := db.Query("SELECT Name,Password FROM users")
	if err != nil {
		log.Fatal(err)
	}
	var name string
	var password string
	for rows.Next() {
		rows.Scan(&name, &password)
		if name == Name {
			if bcrypt.CompareHashAndPassword([]byte(password), []byte(Password)) == nil {
				Name1 = name
				return true, "OK"
			} else {
				return false, "Nepravilnyi porol"
			}
		} else {
			return false, "User does not exist"
		}
	}
	return true, "OK"
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
