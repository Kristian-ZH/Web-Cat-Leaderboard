package data

//User contains data for a submission sender from the database
type User struct {
	ID        int64
	FirstName string
	LastName  string
}

//NewUser is a contructor for User
func NewUser(id int64, firstName, lastName string) *User {
	user := new(User)
	user.ID = id
	user.FirstName = firstName
	user.LastName = lastName
	return user
}
