package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func seedUserData(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var users []User

	json.Unmarshal(byteValue, &users)

	db, err := sql.Open("sqlite3", "data/users.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dropTableQuery := "DROP TABLE IF EXISTS user;"
	_, err = db.Exec(dropTableQuery)
	if err != nil {
		panic(err)
	}

	createTableQuery := "CREATE TABLE IF NOT EXISTS user(id INTEGER PRIMARY KEY, first_name TEXT, last_name TEXT, company TEXT, title TEXT, email TEXT, phone TEXT, dob TEXT, ssn TEXT, salary NUMERIC, admin BOOLEAN);"
	_, err = db.Exec(createTableQuery)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(users); i++ {
		user := users[i]
		insertUserQuery := fmt.Sprintf("INSERT INTO user VALUES (%d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s');", i+1, user.FirstName, user.LastName, user.Company, user.Title, user.Email, user.Phone, user.DOB, user.SSN, user.Salary, user.Admin)
		_, err = db.Exec(insertUserQuery)
		if err != nil {
			panic(err)
		}
	}
}

func getUsers(search string) []User {
	db, err := sql.Open("sqlite3", "data/users.db")
	if err != nil {
		panic(err)
	}
	selectUsersQuery := "SELECT * FROM user WHERE (first_name LIKE '%" + search + "%' OR last_name LIKE '%" + search + "%') AND admin == 'false';"
	rows, err := db.Query(selectUsersQuery)
	if err != nil {
		panic(err)
	}
	var results []User

	for rows.Next() {
		var user User
		_ = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Company, &user.Title, &user.Email, &user.Phone, &user.DOB, &user.SSN, &user.Salary, &user.Admin)
		results = append(results, user)
	}
	rows.Close()
	return results
}

// @felickz testing
func getUser(search string) (User, error) {
	var user User
	db, err := sql.Open("sqlite3", "data/users.db")
	if err != nil {
		return user, err
	}
	defer db.Close()

	selectUserQuery := "SELECT * FROM user WHERE (first_name LIKE ? OR last_name LIKE ?) AND admin == 'false' LIMIT 1;"
	row := db.QueryRow(selectUserQuery, "%"+search+"%", "%"+search+"%")
	err = row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Company, &user.Title, &user.Email, &user.Phone, &user.DOB, &user.SSN, &user.Salary, &user.Admin)
	if err != nil {
		return user, err
	}

	return user, nil
}
