package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/set", set_cookie)
	http.HandleFunc("/get", get_cookie)
	http.HandleFunc("/students", get_students)
	http.HandleFunc("/reg", register)
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
