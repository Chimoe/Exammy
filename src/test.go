package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Accept no argument, this function will read the rcs id stored in the cookie
// Return all available courses of a student user in JSON type
func getStudentCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://ec2-35-153-68-95.compute-1.amazonaws.com")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// check if the user has a valid cookie
	// if not return error message
	if !getCookie(w, r) {
		http.Error(w, "Login again please", http.StatusBadRequest)
		return
	}

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
}

/*
accept a course id (string)
Return all available tests of the course in JSON types
*/
func getStudentTests(w http.ResponseWriter, r *http.Request) {
	if !scrutinize(w, r) {
		return
	}

	// read courseID from request body
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

}

/*
accept a testID (string)
return all questions of the test/quiz in JSON type
*/
func getTestQuestions(w http.ResponseWriter, r *http.Request) {
	if !scrutinize(w, r) {
		return
	}

	// read testID from request body
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

}

/*
accept a JSON containing student's answers and test information
JSON:
{
	"RcsID": "zhongk",
	"TestID": 1,
	"Answers": ["A", "B"]
}
then submit student's answers to database
*/
func submitAnswers(w http.ResponseWriter, r *http.Request) {
	if !scrutinize(w, r) {
		return
	}

	// create a Answer struct
	a := Answer{}

	nameCookie, _ := r.Cookie("username")
	a.RcsID = nameCookie.Value

	// decode contents in r and transfer it to Answer struct
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := Test{}
	t.ID = a.TestID
	inTime, err := t.inTime()

	// check the submission is on time
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if inTime == false {
		http.Error(w, "Time limit exceeded", http.StatusBadRequest)
		return
	}

	// insert the new answer sheet into database
	err = a.submit()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

/*
accept a string containing testID
this function will read the rcs id from the cookie
it will compare the the students' answer with the correct answer
and return the string "{student score}/{total score}"
*/
func autogradeAnswer(w http.ResponseWriter, r *http.Request) {
	if !scrutinize(w, r) {
		return
	}

	nameCookie, _ := r.Cookie("username")
	username := nameCookie.Value

	// read r and transfer the contents in it to int type
	t, _ := ioutil.ReadAll(r.Body)
	ts := string(t)
	testID, _ := strconv.Atoi(ts)

	test := Test{}
	test.ID = testID

	// get the total score of a test
	total, errMsg, errCode := test.totalScore()

	if total == 0 {
		http.Error(w, errMsg, errCode)
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

	//return the string "score/total score"
	var score int
	stmt.QueryRow(testID, username).Scan(&score)
	fmt.Fprint(w, score)
	fmt.Fprint(w, "/")
	fmt.Fprint(w, total)
	w.WriteHeader(http.StatusOK)
}
