package domain

type Error struct {
	Message string `json:"message"`
}

func NewError(m string) Error {
	return Error{Message: m}
}
