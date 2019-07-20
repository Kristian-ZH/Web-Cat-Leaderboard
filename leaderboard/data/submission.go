package data

//Submission contains data for a submission from the database
type Submission struct {
	User       User
	Assignment Assignment
	Score      float64
}

//NewSubmission is a contructor for Submission
func NewSubmission(user User, assignment Assignment, score float64) *Submission {
	submission := new(Submission)
	submission.User = user
	submission.Assignment = assignment
	submission.Score = score
	return submission
}
