package main

import (
	"github.com/satori/go.uuid"
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
func setCookie(w http.ResponseWriter, r *http.Request, RcsID string) {
	if cookies == nil {
		cookies = make(map[string]string)
	}
	expiration := time.Now().Add(time.Hour)
	u1 := uuid.Must(uuid.NewV4()).String()
	nameCookie := http.Cookie{Name: "username", Value: string(RcsID), Expires: expiration}
	cookie := http.Cookie{Name: string(RcsID), Value: u1, Expires: expiration}
	cookies[cookie.Name] = cookie.Value
	http.SetCookie(w, &nameCookie)
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
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
		return false
	} else {
		cookie, _ := r.Cookie(nameCookie.Value)
		if cookie.Value != cookies[cookie.Name] {
			return false
		}
		return true
	}
}
