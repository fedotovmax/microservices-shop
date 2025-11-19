package domain

import "regexp"

type CreateUser struct {
	Email string
}

type User struct {
	ID    string
	Email string
}

type UserCreatedEvent struct {
	ID string
}

var PhoneRegexp = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)
