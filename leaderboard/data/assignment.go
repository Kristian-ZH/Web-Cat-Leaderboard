package data

//Assignment contains data for an assignment from the database
type Assignment struct {
	Name string
}

//NewAssignment is a contructor for Assignment
func NewAssignment(name string) *Assignment {
	assignment := new(Assignment)
	assignment.Name = name
	return assignment
}
