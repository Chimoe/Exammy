package main

import (
	"log"
	"net/http"
)

func main() {
	//http.HandleFunc("/set", setCookie)
	//http.HandleFunc("/check", print_all_cookies)

	http.HandleFunc("/reg", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/courses", getStudentCourses)
	http.HandleFunc("/tests", getStudentTests)
	http.HandleFunc("/questions", getTestQuestions)
	http.HandleFunc("/answers", submitAnswers)
	http.HandleFunc("/grade", autogradeAnswer)
	http.HandleFunc("/NM$L", secret)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
