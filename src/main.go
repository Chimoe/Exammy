package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/set", setCookie)
	//http.HandleFunc("/get", getCookie)
	http.HandleFunc("/students", get_student)
	http.HandleFunc("/reg", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/courses", getStudentCourses)
	http.HandleFunc("/tests", getStudentTests)
	//http.HandleFunc("/check", print_all_cookies)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
