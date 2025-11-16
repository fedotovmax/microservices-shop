package domain

type CreateUser struct {
	Email     string
	FirstName string
	LastName  string
}

type User struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
}
