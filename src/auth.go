package main

import (
	"github.com/satori/go.uuid"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	cookies map[string]string
}

var s *Session
var once sync.Once

/*
Get the valid cookies list
*/
func getSession() *Session {
	once.Do(func() {
		s = &Session{cookies: make(map[string]string)}
	})
	return s
}

/*
Add a new valid cookie into cookies list
*/
func (s Session) addCookie(cookie http.Cookie) {
	s.cookies[cookie.Name] = cookie.Value
}

/*
Check if the user's cookie is valid
*/
func (s Session) checkCookie(cookie http.Cookie) bool {
	return cookie.Value == s.cookies[cookie.Name]
}

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
	session := getSession()
	expiration := time.Now().Add(time.Hour)
	u1 := uuid.Must(uuid.NewV4()).String()
	nameCookie := http.Cookie{Name: "username", Value: string(RcsID), Expires: expiration}
	cookie := http.Cookie{Name: string(RcsID), Value: u1, Expires: expiration}
	session.addCookie(cookie)
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
	session := getSession()
	nameCookie, _ := r.Cookie("username")
	if nameCookie == nil {
		return false
	} else {
		cookie, _ := r.Cookie(nameCookie.Value)
		return session.checkCookie(*cookie)
	}
}
