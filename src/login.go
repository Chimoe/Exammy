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

func login(w http.ResponseWriter, r *http.Request) {
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

	db, err := sql.Open("mysql",
		dataSourceName)
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
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "WRONG RCS ID")
		return
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(s.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "WRONG PASSWORD")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func register(w http.ResponseWriter, r *http.Request) {
	s := Student{}
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// pwd, _ := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	// s.Password = string(pwd)
	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
