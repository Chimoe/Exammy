package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Student struct {
	RcsID     string
	FirstName string
	LastName  string
	Password  string
}

/*  a login API
accept a JSON containing student's rcs id and password
return plain text indicating if the id and password are correct
*/
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}

	s := Student{}
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if s.RcsID == "" || s.Password == "" {
		http.Error(w, "Please send a rcs_id and a password", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT password FROM student WHERE rcs_id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var hashedPassword []byte
	stmt.QueryRow(s.RcsID).Scan(&hashedPassword)

	if len(hashedPassword) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "WRONG RCS ID")
		return
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(s.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "WRONG PASSWORD")
		return
	}

	setCookie(w, r)
	fmt.Fprint(w, "OK")
}

/*  a register API
accept a JSON containing student's rcs id, first name, last name and password
create a new student account and insert it into database
return plain text indicating if the account are registered correctly
*/
func register(w http.ResponseWriter, r *http.Request) {
	s := Student{}
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pwd, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.Password = string(pwd)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT student SET rcs_id = ?, first_name = ?, last_name = ?, password = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(s.RcsID, s.FirstName, s.LastName, s.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	LastInsertId, err := res.LastInsertId()
	if err == nil {
		fmt.Println("LastInsertId:", LastInsertId)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
	/*
		w.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	*/
}
