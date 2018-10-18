package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/set", set_cookie)
	//http.HandleFunc("/get", get_cookie)
	http.HandleFunc("/students", get_student)
	http.HandleFunc("/reg", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/courses", get_student_courses)
	//http.HandleFunc("/check", print_all_cookies)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
