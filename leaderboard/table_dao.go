package leaderboard

import (
	"fmt"
	"os"

	"github.com/leaderboard/Web-Cat-Leaderboard/leaderboard/data"

	"github.com/leaderboard/Web-Cat-Leaderboard/database"
)

var (
	courseName = os.Getenv("COURSE_NAME") //COURSE_NAME=MJT2019
)

//GetUserNames returns slice with all users from the current Course
func GetUserNames() []data.User {
	db := database.DbConn()
	selectQuery := fmt.Sprintf("SELECT user.OID, user.CFIRSTNAME, user.CLASTNAME "+
		"FROM TUSER user "+
		"INNER JOIN TCOURSESTUDENT cu ON user.OID=cu.CID1 "+
		"INNER JOIN TCOURSEOFFERING coff ON cu.CID=coff.OID "+
		"INNER JOIN TCOURSE course ON coff.CCOURSEID=course.OID "+
		"WHERE CACCESSLEVEL='0' AND coff.CCRN='%s';", courseName)
	selDB, err := db.Query(selectQuery)
	if err != nil {
		panic(err.Error())
	}
	userNames := []data.User{}
	for selDB.Next() {
		var id int64
		var firstName, lastName string
		err = selDB.Scan(&id, &firstName, &lastName)
		if err != nil {
			panic(err.Error())
		}
		user := data.NewUser(id, firstName, lastName)
		userNames = append(userNames, *user)
	}
	defer db.Close()
	return userNames
}

//GetAssignments returns slice with all assignments from the current Course
func GetAssignments() []data.Assignment {
	db := database.DbConn()
	selectQuery := fmt.Sprintf("SELECT ass.CASSIGNMENTNAME "+
		"FROM TASSIGNMENT ass "+
		"INNER JOIN TASSIGNMENTOFFERING assoff ON ass.OID=assoff.CASSIGNMENTID "+
		"INNER JOIN TCOURSEOFFERING coff ON assoff.CCOURSEOFFERINGID=coff.OID "+
		"WHERE CASSIGNMENTURL IS NOT NULL AND coff.CCRN='%s'"+
		"ORDER BY coff.OID;", courseName)
	selDB, err := db.Query(selectQuery)
	if err != nil {
		panic(err.Error())
	}
	assignments := []data.Assignment{}
	for selDB.Next() {
		var name string
		err = selDB.Scan(&name)
		if err != nil {
			panic(err.Error())
		}
		assignment := data.NewAssignment(name)
		assignments = append(assignments, *assignment)
	}
	defer db.Close()
	return assignments
}

//GetSubmissions returns slice with all submissions from the current Course
func GetSubmissions() []data.Submission {
	db := database.DbConn()
	selectQuery := fmt.Sprintf("SELECT user.OID, user.CFIRSTNAME, user.CLASTNAME ,ass.CASSIGNMENTNAME, ( COALESCE(subresult.CCORRECTNESSSCORE, 0) + COALESCE(subresult.CTOOLSCORE, 0) + COALESCE(subresult.CTASCORE ,0)) AS score "+
		"FROM TSUBMISSION sub "+
		"INNER JOIN TSUBMISSIONRESULT subresult ON subresult.OID=sub.CRESULTID "+
		"INNER JOIN TASSIGNMENT ass ON ass.OID=sub.CASSIGNMENTID "+
		"INNER JOIN TUSER user ON user.OID=sub.CUSERID "+
		"INNER JOIN TCOURSESTUDENT cu ON user.OID=cu.CID1 "+
		"INNER JOIN TCOURSEOFFERING coff ON cu.CID=coff.OID "+
		"INNER JOIN TCOURSE course ON coff.CCOURSEID=course.OID "+
		"WHERE subresult.CISMOSTRECENT='1' AND ass.CASSIGNMENTURL IS NOT NULL AND user.CACCESSLEVEL='0' AND coff.CCRN='%s';", courseName)
	selDB, err := db.Query(selectQuery)
	if err != nil {
		panic(err.Error())
	}

	var uid int64
	var firstName, lastName, assignmentName string
	var score float64

	submissions := []data.Submission{}
	for selDB.Next() {
		err = selDB.Scan(&uid, &firstName, &lastName, &assignmentName, &score)
		if err != nil {
			panic(err.Error())
		}
		user := data.NewUser(uid, firstName, lastName)
		assignment := data.NewAssignment(assignmentName)
		submission := data.NewSubmission(*user, *assignment, score)
		submissions = append(submissions, *submission)
	}
	defer db.Close()
	return submissions
}
