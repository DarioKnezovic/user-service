package user

// User represents a user entity.
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// UserService represents the user service interface.
type UserService interface {
	RegisterUser(newUser User) (*User, error)
}
