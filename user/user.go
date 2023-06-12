package user

// User represents a registered user
type User struct {
	Username string
	Password string
}

// UserMap stores user information
type UserMap map[string]*User

// NewUserMap creates a new instance of UserMap
func NewUserMap() UserMap {
	return make(UserMap)
}

// AddUser adds a user to the UserMap
func (m UserMap) AddUser(username, password string) {
	m[username] = &User{
		Username: username,
		Password: password,
	}
}
