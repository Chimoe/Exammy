package main

import (
	"fmt"
	"net/http"
	"time"
)

func set_cookie(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(time.Hour)
	cookie := http.Cookie{Name: "username", Value: "asdasdasd", Expires: expiration}
	http.SetCookie(w, &cookie)
}

func get_cookie(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("username")
	if cookie == nil {
		fmt.Fprint(w, "NM$L")
		return
	}
	fmt.Fprint(w, cookie)
}

func get_students(w http.ResponseWriter, r *http.Request) {
	/* student1 := Student{"ziyi lu", "", false}
	student2 := Student{"yanlin zhu", "", false}
	student3 := Student{"jingfei zhou", "", false}
	students := []Student{student1, student2}
	students = append(students, student3)
	js, err := json.Marshal(students)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js) */
}
