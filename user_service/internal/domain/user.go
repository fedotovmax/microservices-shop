package domain

type CreateUser struct {
	Email     string
	FirstName string
	LastName  string
}

type UserCreatedEvent struct {
	ID string
}
