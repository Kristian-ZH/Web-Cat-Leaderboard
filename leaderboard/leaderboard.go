package leaderboard

import (
	"strconv"

	"github.com/leaderboard/Web-Cat-Leaderboard/leaderboard/data"
)

var (
	userNames   []data.User
	assignments []data.Assignment
	table       [][]string
)

const (
	//ScoreTableOffseet is the offset of the Total scores in the table
	ScoreTableOffseet = 2
	//HeadersTableOffset is the offset of the headers in the table
	HeadersTableOffset = 1
)

//GetLeaderboardTable retuns string matrix with the scores from all submissions
func GetLeaderboardTable() *[][]string {
	assignments = GetAssignments()
	userNames = GetUserNames()
	table = make([][]string, len(userNames)+1)
	for i := range table {
		table[i] = make([]string, len(assignments)+ScoreTableOffseet)
	}

	table[0][1] = "Total"
	for idx, assignment := range assignments {
		table[0][idx+ScoreTableOffseet] = assignment.Name
	}
	for idx, user := range userNames {
		table[idx+1][0] = user.FirstName + " " + user.LastName
	}

	submissions := GetSubmissions()
	for _, submission := range submissions {
		saveSubmissionInTable(&submission)
	}

	applyTotalScores()
	sortScores()
	return &table
}

func saveSubmissionInTable(submission *data.Submission) {
	i := findUserIndex(&submission.User)
	j := findAssignmentIndex(&submission.Assignment)
	table[i][j] = strconv.FormatFloat(submission.Score, 'f', 1, 64)
}

func findUserIndex(user *data.User) int {
	userIndex := 0
	for i := HeadersTableOffset; i < len(userNames)+HeadersTableOffset; i++ {
		if table[i][0] == user.FirstName+" "+user.LastName {
			userIndex = i
			break
		}
	}

	return userIndex
}

func findAssignmentIndex(assignment *data.Assignment) int {
	assignmentIndex := 0
	for j := HeadersTableOffset; j < len(assignments)+ScoreTableOffseet; j++ {
		if table[0][j] == assignment.Name {
			assignmentIndex = j
			break
		}
	}

	return assignmentIndex
}

func applyTotalScores() {
	for i := HeadersTableOffset; i < len(userNames)+HeadersTableOffset; i++ {
		var userTotal float64
		for j := ScoreTableOffseet; j < len(assignments)+ScoreTableOffseet; j++ {
			if score, err := strconv.ParseFloat(table[i][j], 64); err == nil {
				userTotal += score
			}
		}
		table[i][HeadersTableOffset] = strconv.FormatFloat(userTotal, 'f', 1, 64)
	}
}

func sortScores() {
	scoreIndex := 1
	for i := HeadersTableOffset + 1; i < len(userNames)+HeadersTableOffset; i++ {
		j := i
		for j > 1 {
			currentScore, _ := strconv.ParseFloat(table[j][scoreIndex], 64)
			previousScore, _ := strconv.ParseFloat(table[j-1][scoreIndex], 64)
			if previousScore < currentScore {
				table[j-1], table[j] = table[j], table[j-1]
			}
			j--
		}
	}
}
