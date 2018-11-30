package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"strconv"
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

func (t Test) totalScore(w http.ResponseWriter) int {
	// open database
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 0
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT sum(q.score)
										FROM question q
										WHERE q.testID = ?`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 0
	}

	defer stmt.Close()

	var totalScore int
	stmt.QueryRow(t.ID).Scan(&totalScore)
	if totalScore == 0 {
		http.Error(w, "Test does not exist", http.StatusBadRequest)
	}
	return totalScore
}

type Question struct {
	TestID      int
	QuestionNum int
	Text        string
	Answer      string
	Score       int
}

type Answer struct {
	RcsID   string
	TestID  int
	Answers []string
}

/*
Return all available courses of a student user
*/
func getStudentCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://ec2-35-153-68-95.compute-1.amazonaws.com")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// check if the user has a valid cookie
	// if not return error message
	if getCookie(w, r) {
		// get user id from cookie
		nameCookie, _ := r.Cookie("username")
		username := nameCookie.Value

		// open database
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// get the course information
		stmt, err := db.Prepare(`select c.course_id, c.subject, c.number, c.name
		from student_course sc join user s on s.student_id = sc.student_id
		join course c on c.course_id = sc.course_id
		where s.rcs_id = ?`)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		// input the username
		rows, _ := stmt.Query(username)
		defer rows.Close()

		// read the course information and store it in a list of Course struct
		var courses []Course
		// read a course in each iteration
		for rows.Next() {
			var c Course
			rows.Scan(&c.ID, &c.Subject, &c.Number, &c.Name)
			courses = append(courses, c)
		}

		w.Header().Set("Content-Type", "application/json")

		// marshal courses to JSON data
		js, err := json.Marshal(courses)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		// send JSON data to front-end
		w.Write(js)
	} else {
		fmt.Fprint(w, r.Cookies())
		http.Error(w, "Login again please", http.StatusBadRequest)
	}
}

/*
accept a course id
Return all available tests of the course
*/
func getStudentTests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://ec2-35-153-68-95.compute-1.amazonaws.com")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// check if the user has a valid cookie
	// if not return error message
	if getCookie(w, r) {
		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			return
		}
		// read courseID from r
		courseID, _ := ioutil.ReadAll(r.Body)

		// open database
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// get the test information
		stmt, err := db.Prepare(`SELECT * FROM test WHERE course_id = ? ORDER BY test.name`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		// input the courseID
		rows, _ := stmt.Query(courseID)
		defer rows.Close()

		// read the test information and store it in a list of Test struct
		var tests []Test
		// read a test in each iteration
		for rows.Next() {
			var t Test
			rows.Scan(&t.ID, &t.CourseID, &t.Name, &t.StartT, &t.EndT)
			tests = append(tests, t)
		}

		w.Header().Set("Content-Type", "application/json")
		// marshal tests to JSON data
		js, err := json.Marshal(tests)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		// send JSON data to front-end
		w.Write(js)

	} else {
		http.Error(w, "Login again please", http.StatusBadRequest)
	}
}

/*
accept a testID
return all questions of the test/quiz
*/
func getTestQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://ec2-35-153-68-95.compute-1.amazonaws.com")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// check if the user has a valid cookie
	// if not return error message
	if getCookie(w, r) {
		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			return
		}
		// read testID from r
		testID, _ := ioutil.ReadAll(r.Body)

		// open database
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// get the question information
		stmt, err := db.Prepare(`SELECT q.questionNum, q.text, q.score
								FROM test t JOIN question q ON q.testID = t.test_id
								WHERE t.start_t < CONVERT_TZ(NOW(),'UTC','America/New_York')
									AND t.test_id = ?
								ORDER BY q.questionNum`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		// input the testID
		rows, _ := stmt.Query(testID)
		defer rows.Close()

		// read the question information and store it in a list of Question struct
		var questions []Question
		// read a question in each iteration
		for rows.Next() {
			var q Question
			rows.Scan(&q.QuestionNum, &q.Text, &q.Score)
			questions = append(questions, q)
		}

		w.Header().Set("Content-Type", "application/json")
		// marshal questions to JSON data
		js, err := json.Marshal(questions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		// send JSON data to front-end
		w.Write(js)

	} else {
		http.Error(w, "Login again please", http.StatusBadRequest)
	}
}

/*
accept a JSON containing student's answers and test information
then submit student's answers with his/her id to database
*/
func submitAnswers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://ec2-35-153-68-95.compute-1.amazonaws.com")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// check if the user has a valid cookie
	// if not return error message
	if getCookie(w, r) {
		// create a Answer struct
		a := Answer{}

		nameCookie, _ := r.Cookie("username")
		a.RcsID = nameCookie.Value

		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			return
		}
		// decode contents in r and transfer it to Answer struct
		err := json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// open database
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// get the test information
		stmt, err := db.Prepare(`SELECT count(*)
						 FROM test t
						 WHERE t.start_t < CONVERT_TZ(NOW(),'UTC','America/New_York')
 						 AND t.end_t > CONVERT_TZ(NOW(),'UTC','America/New_York')
						 AND t.test_id = ?`)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer stmt.Close()

		// read the test expiration time
		var inTime int
		stmt.QueryRow(a.TestID).Scan(&inTime)
		// check if the test limit exceeded
		if inTime == 0 {
			http.Error(w, "Time Limit Exceeded", http.StatusBadRequest)
			return
		}

		// insert a new answer sheet into database
		stmt, err = db.Prepare(`INSERT answer_history SET rcs_id = ?, testID = ?, questionNum = ?,
                          					answer = ?`)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer stmt.Close()

		// insert student id, test id, question number, and corresponding answer to database
		for i, ans := range a.Answers {
			_, err := stmt.Exec(a.RcsID, a.TestID, i+1, ans)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Login again please", http.StatusBadRequest)
	}
}

/*
accept a string containing testID
then submit student's grade of the given test
*/
func autogradeAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://ec2-35-153-68-95.compute-1.amazonaws.com")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// check if the user has a valid cookie
	// if not return error message
	if getCookie(w, r) {
		nameCookie, _ := r.Cookie("username")
		username := nameCookie.Value

		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			return
		}
		// read r and transfer the contents in it to int type
		t, _ := ioutil.ReadAll(r.Body)
		ts := string(t)
		testID, _ := strconv.Atoi(ts)

		test := Test{}
		test.ID = testID

		totalScoreV := test.totalScore(w)
		if totalScoreV == 0 {
			return
		}

		// check the student's answers

		// open database
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		stmt, err := db.Prepare(`SELECT sum(q.score)
										FROM question q JOIN answer_history ah on q.testID = ah.testID
                                        and q.questionNum = ah.questionNum
                                        and q.answer = ah.answer
										WHERE q.testID = ? and ah.rcs_id = ?`)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var score int
		stmt.QueryRow(testID, username).Scan(&score)
		fmt.Fprint(w, score)
		fmt.Fprint(w, "/")
		fmt.Fprint(w, totalScoreV)
		w.WriteHeader(http.StatusOK)

	} else {
		http.Error(w, "Login again please", http.StatusBadRequest)
	}
}
