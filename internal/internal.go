package internal

func ConfirmSignup(Name string, Email string, Password string, RewrittenPassword string) (bool, string) {
	// Надо дописать всякие условия
	if RewrittenPassword != Password {
		return false, "Passwords don't match, write again."
	}
	if len(Name) < 3 || len(Password) < 7 {
		return false, "Name/Password doesn't have enough characters. Minimum for name is 3 and for password is 7."
	}
	// db, err := sql.Open("sqlite3", "./example.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// tableName := "users"

	// query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", tableName)
	// row1 := db.QueryRow(query)

	// var name string
	// err = row1.Scan(&name)
	// if err == sql.ErrNoRows {
	// } else if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	rows, err := db.Query("SELECT Name, Email FROM users")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	var name string
	// 	var email string
	// 	for rows.Next() {
	// 		rows.Scan(&name, &email)
	// 		fmt.Println(name, email)
	// 		if name == Name {
	// 			db.Close()
	// 			return false, "That name is already being used"
	// 		} else if Email == email {
	// 			return false, "That Email is already being used"
	// 		}
	// 	}
	// }

	return true, "OK"
}

func ConfirmSignin(Name string, Password string) (bool, string) {
	// Надо дописать всякие условия
	return true, "OK"
}
