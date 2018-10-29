package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"time"
)

var cookies map[string]string

/* a helper function
print all username and uuid stored in cookies
*/
/*
func print_all_cookies(w http.ResponseWriter, r *http.Request) {
	for username := range cookies {
		fmt.Fprint(w, "username: ", username, " ", "uuid: ", cookies[username], "\n")
	}
}
*/

/*
accept a string containing username
create a uuid pairing with the input username and store them into global map cookies
*/
func setCookie(w http.ResponseWriter, r *http.Request) {
	if cookies == nil {
		cookies = make(map[string]string)
	}
	username, _ := ioutil.ReadAll(r.Body)
	expiration := time.Now().Add(time.Hour)
	u1 := uuid.Must(uuid.NewV4()).String()
	nameCookie := http.Cookie{Name: "username", Value: string(username), Expires: expiration}
	cookie := http.Cookie{Name: string(username), Value: u1, Expires: expiration}
	cookies[cookie.Name] = cookie.Value
	http.SetCookie(w, &nameCookie)
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

/*
accept two cookies
nameCookie: the first one has Name: "username" and Value: string(username)
Cookie: the second one has Name: string(username) and Value: uuid
then check if the cookie we get has the valid uuid
*/
func getCookie(w http.ResponseWriter, r *http.Request) bool {
	if cookies == nil {
		cookies = make(map[string]string)
	}
	nameCookie, _ := r.Cookie("username")
	if nameCookie == nil {
		// fmt.Fprint(w, "NM$L")
		return false
	} else {
		cookie, _ := r.Cookie(nameCookie.Value)
		if cookie.Value != cookies[cookie.Name] {
			// w.WriteHeader(http.StatusUnauthorized)
			// fmt.Fprint(w, "INVALID COOKIE")
			return false
		}
		// w.WriteHeader(http.StatusOK)
		// fmt.Fprint(w, "OK")
		return true
	}
}

func get_student(w http.ResponseWriter, r *http.Request) {
	if getCookie(w, r) {
		fmt.Fprint(w, "OK")
	} else {
		fmt.Fprint(w, "NM$L")
	}
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
