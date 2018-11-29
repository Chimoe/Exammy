package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type User struct {
	RcsID     string
	FirstName string
	LastName  string
	Password  string
	Identity  bool //true: student   false: instructor
}

/*
Encrypt the user password
*/
func (u User) encryptedPassword() {
	pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	u.Password = string(pwd)
}

/*
Check if the user enters ID and password
*/
func (u User) empty() bool {
	return u.RcsID == "" || u.Password == ""
}

/*  a login API
accept a JSON containing student's rcs id and password
return plain text indicating if the id and password are correct
*/
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	//var a []byte
	// var a []byte

	// r.Body.Read(a)
	// fmt.Print(string(a))
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	// create a User struct
	u := User{}
	// decode contents in r and transfer it to User struct
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if u.empty() {
		http.Error(w, "Please send a rcs_id and a password", http.StatusBadRequest)
		return
	}

	// open database
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// get the user information
	stmt, err := db.Prepare("SELECT password FROM user WHERE rcs_id = ? AND identity = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var hashedPassword []byte
	stmt.QueryRow(u.RcsID, u.Identity).Scan(&hashedPassword)

	// check if the user id is valid
	if len(hashedPassword) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "WRONG RCS ID")
		return
	}

	// check if the user password is valid
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(u.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "WRONG PASSWORD")
		return
	}

	// set a cookie for the valid user account
	setCookie(w, r, u.RcsID)
	fmt.Fprint(w, "OK")
}

/*  a register API
accept a JSON containing student's rcs id, first name, last name and password
create a new student account and insert it into database
return plain text indicating if the account are registered correctly
*/
func register(w http.ResponseWriter, r *http.Request) {
	// create a User struct
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	u := User{}
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	// decode contents in r and transfer it to User struct
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// encrypt user password
	u.encryptedPassword()

	// open database
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// insert a new user account into database
	stmt, err := db.Prepare(`INSERT user SET rcs_id = ?, first_name = ?, last_name = ?, password = ?
								   , identity = ?`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.RcsID, u.FirstName, u.LastName, u.Password, u.Identity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/* LastInsertId, err := res.LastInsertId()
	 if err == nil {
		fmt.Println("LastInsertId:", LastInsertId)
	} */

	w.WriteHeader(http.StatusOK)
}
