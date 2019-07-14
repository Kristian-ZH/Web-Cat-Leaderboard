package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type SubmissionResult struct {
	UserName   string
	Assignment string
	score      float64
}

type User struct {
	ID        int64
	FirstName string
	LastName  string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := ""
	dbPass := ""
	dbName := ""
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	db.SetMaxIdleConns(0)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	assignmentNames := getAssignments()
	userNames := getUserNames()
	table := make([][]string, len(userNames)+1)
	for i := range table {
		table[i] = make([]string, len(assignmentNames)+1)
	}

	for idx, assignment := range assignmentNames {
		table[0][idx+1] = assignment
	}
	for idx, user := range userNames {
		table[idx+1][0] = user.FirstName + " " + user.LastName
	}

	db := dbConn()
	selectQuery := "SELECT user.OID, user.CFIRSTNAME, user.CLASTNAME ,ass.CASSIGNMENTNAME, ( COALESCE(subresult.CCORRECTNESSSCORE, 0) + COALESCE(subresult.CTOOLSCORE, 0) + COALESCE(subresult.CTASCORE ,0)) AS score " +
		"FROM TSUBMISSION sub " +
		"INNER JOIN TSUBMISSIONRESULT subresult ON subresult.OID=sub.CRESULTID " +
		"INNER JOIN TASSIGNMENT ass ON ass.OID=sub.CASSIGNMENTID " +
		"INNER JOIN TUSER user ON user.OID=sub.CUSERID " +
		"WHERE subresult.CISMOSTRECENT='1';"
	selDB, err := db.Query(selectQuery)
	if err != nil {
		panic(err.Error())
	}

	var uid int64
	var firstName, lastName, assignmentName string
	var score float64

	for selDB.Next() {
		err = selDB.Scan(&uid, &firstName, &lastName, &assignmentName, &score)
		if err != nil {
			break
		}

		user := User{uid, firstName, lastName}
		saveSubmissionInTable(&user, assignmentName, score, &table, len(userNames), len(assignmentNames))
	}

	tmpl.ExecuteTemplate(w, "Index", table)
	defer db.Close()
}

func getUserNames() []User {
	db := dbConn()
	selDB, err := db.Query("SELECT OID, CFIRSTNAME, CLASTNAME FROM TUSER;")
	if err != nil {
		panic(err.Error())
	}
	user := User{}
	userNames := []User{}
	for selDB.Next() {
		var id int64
		var firstName, lastName string
		err = selDB.Scan(&id, &firstName, &lastName)
		if err != nil {
			panic(err.Error())
		}
		user.ID = id
		user.FirstName = firstName
		user.LastName = lastName
		userNames = append(userNames, user)
	}
	defer db.Close()
	return userNames
}

func getAssignments() []string {
	db := dbConn()
	selDB, err := db.Query("SELECT CASSIGNMENTNAME FROM TASSIGNMENT ORDER BY OID;")
	if err != nil {
		panic(err.Error())
	}
	assignmentNames := []string{}
	for selDB.Next() {
		var name string
		err = selDB.Scan(&name)
		if err != nil {
			panic(err.Error())
		}
		assignmentNames = append(assignmentNames, name)
	}
	defer db.Close()
	return assignmentNames
}

func saveSubmissionInTable(user *User, assignmentName string, score float64, table *[][]string, usersNum, assignmentsNum int) {
	i := findUserIndex(user, table, usersNum)
	j := findAssignmentIndex(assignmentName, table, assignmentsNum)
	(*table)[i][j] = strconv.FormatFloat(score, 'f', 1, 64)
	fmt.Printf("%s_%s  %s  %f :  i=%d : j=%d\n", user.FirstName, user.LastName, assignmentName, score, i, j)
}

func findUserIndex(user *User, table *[][]string, usersNum int) int {
	userIndex := 0
	for i := 1; i <= usersNum; i++ {
		if (*table)[i][0] == user.FirstName+" "+user.LastName {
			userIndex = i
			break
		}
	}

	return userIndex
}

func findAssignmentIndex(assignmentName string, table *[][]string, assignmentsNum int) int {
	assignmentIndex := 0
	for j := 1; j <= assignmentsNum; j++ {
		if (*table)[0][j] == assignmentName {
			assignmentIndex = j
			break
		}
	}

	return assignmentIndex
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.ListenAndServe(":8080", nil)
}
