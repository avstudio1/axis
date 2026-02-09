package workspace

type User struct {
	Email string
	Name  string
}

type Service interface {
	GetUser(email string) (*User, error)
}
