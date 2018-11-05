package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"time"
)

type Course struct {
	ID      int
	Subject string
	Number  string
	Name    string
}

type Test struct {
	ID       int
	CourseID int
	Name     string
	StartT   time.Time
	EndT     time.Time
}

type Question struct {
	TestID      int
	QuestionNum int
	Text        string
	Answer      string
	Score       int
}

func getStudentCourses(w http.ResponseWriter, r *http.Request) {
	if getCookie(w, r) {
		nameCookie, _ := r.Cookie("username")
		username := nameCookie.Value

		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		stmt, err := db.Prepare(`select c.course_id, c.subject, c.number, c.name
		from student_course sc join user s on s.student_id = sc.student_id
		join course c on c.course_id = sc.course_id
		where s.rcs_id = ?`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		rows, _ := stmt.Query(username)
		defer rows.Close()

		var courses []Course
		for rows.Next() {
			var c Course
			rows.Scan(&c.ID, &c.Subject, &c.Number, &c.Name)
			courses = append(courses, c)
		}

		w.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(courses)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	} else {
		http.Error(w, "Login again please", http.StatusBadRequest)
	}
}

func getStudentTests(w http.ResponseWriter, r *http.Request) {
	if getCookie(w, r) {
		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			return
		}
		courseID, _ := ioutil.ReadAll(r.Body)

		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		stmt, err := db.Prepare(`SELECT * FROM test WHERE course_id = ? ORDER BY test.name`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		rows, _ := stmt.Query(courseID)
		defer rows.Close()

		var tests []Test
		for rows.Next() {
			var t Test
			rows.Scan(&t.ID, &t.CourseID, &t.Name, &t.StartT, &t.EndT)
			tests = append(tests, t)
		}

		w.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(tests)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(js)

	} else {
		http.Error(w, "Login again please", http.StatusBadRequest)
	}
}

func getTestQuestions(w http.ResponseWriter, r *http.Request) {
	if getCookie(w, r) {
		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			return
		}

		testID, _ := ioutil.ReadAll(r.Body)

		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		stmt, err := db.Prepare(`SELECT q.questionNum, q.text, q.score
								FROM test t JOIN question q ON q.testID = t.test_id
								WHERE t.start_t < CONVERT_TZ(NOW(),'UTC','America/New_York')
									AND t.end_t > CONVERT_TZ(NOW(),'UTC','America/New_York')
									AND t.test_id = ?
								ORDER BY q.questionNum`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		rows, _ := stmt.Query(testID)
		defer rows.Close()

		var questions []Question
		for rows.Next() {
			var q Question
			rows.Scan(&q.QuestionNum, &q.Text, &q.Score)
			questions = append(questions, q)
		}

		w.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(questions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(js)

	} else {
		http.Error(w, "Login again please", http.StatusBadRequest)
	}
}
