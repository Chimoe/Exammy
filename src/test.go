package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

type Course struct {
	Subject string
	Number  string
	Name    string
}

func get_student_courses(w http.ResponseWriter, r *http.Request) {
	if get_cookie(w, r) {
		// fmt.Fprint(w, "OK")
		nameCookie, _ := r.Cookie("username")
		username := nameCookie.Value

		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		stmt, err := db.Prepare(`select c.subject, c.number, c.name
		from student_course sc join student s on s.student_id = sc.student_id
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
			rows.Scan(&c.Subject, &c.Number, &c.Name)
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
		http.Error(w, "Please send a request body", http.StatusBadRequest)
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
