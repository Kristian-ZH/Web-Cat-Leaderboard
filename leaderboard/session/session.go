package session

const (
	//SessionCookieName is the name of session cookie
	SessionCookieName = "wosid"
)

//Session contains data needed to check if session is valid
type Session struct {
	SessionID string
	UserID    int64
}

//NewSession is a contructor for Session
func NewSession(sessionID string, userID int64) *Session {
	session := new(Session)
	session.SessionID = sessionID
	session.UserID = userID
	return session
}
