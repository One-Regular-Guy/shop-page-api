package domain

type User struct {
	ID       string
	Name     string
	Username string
	Email    string
	Password string
}

type UserRepository interface {
	IsUsernameTaken(username string) (bool, error)
	IsEmailTaken(email string) (bool, error)
	CreateUser(user *User) error
	CheckPassword(username, password string) (bool, error)
}
