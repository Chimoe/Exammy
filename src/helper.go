package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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

//return the total score of a test
func (t Test) totalScore() (total int, errMsg string, errCode int) {
	// open database
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return 0, err.Error(), http.StatusInternalServerError
	}
	defer db.Close()

	//prepare the SQL statement
	stmt, err := db.Prepare(`SELECT sum(q.score)
										FROM question q
										WHERE q.testID = ?`)

	if err != nil {
		return 0, err.Error(), http.StatusInternalServerError
	}

	defer stmt.Close()

	// execute the statement and return the total score
	var totalScore int
	stmt.QueryRow(t.ID).Scan(&totalScore)
	if totalScore == 0 {
		return 0, "Test does not exist", http.StatusBadRequest
	}
	return totalScore, "", 0
}

//check if the current time is in the time slice of the test
func (t Test) inTime() (start bool, err error) {
	// open database
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return false, err
	}
	defer db.Close()

	//prepare the SQL statement
	stmt, err := db.Prepare(`SELECT count(*)
						 FROM test t
						 WHERE t.start_t < CONVERT_TZ(NOW(),'UTC','America/New_York')
 						 AND t.end_t > CONVERT_TZ(NOW(),'UTC','America/New_York')
						 AND t.test_id = ?`)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	// read the test information
	var inTime int
	stmt.QueryRow(t.ID).Scan(&inTime)
	// check if the test limit exceeded
	return inTime >= 1, nil
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

// submit the answer object to the database
func (a Answer) submit() (err error) {
	// open database
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	defer db.Close()

	// insert a new answer sheet into database
	stmt, err := db.Prepare(`INSERT answer_history SET rcs_id = ?, testID = ?, questionNum = ?,
                          					answer = ?`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	// insert student id, test id, question number, and corresponding answer to database
	for i, ans := range a.Answers {
		_, err := stmt.Exec(a.RcsID, a.TestID, i+1, ans)

		if err != nil {
			return err
		}
	}
	return nil
}

// scrutinize the input, check the cookie and ensure the request body is not empty
func scrutinize(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "http://ec2-35-153-68-95.compute-1.amazonaws.com")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// check if the user has a valid cookie
	// if not return error message
	if !getCookie(w, r) {
		http.Error(w, "Login again please", http.StatusBadRequest)
		return false
	}

	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return false
	}

	return true
}
