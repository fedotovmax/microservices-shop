package domain

const INVALID_BODY = "Тело запроса имеет недопустимый формат!"

type Error struct {
	Message string `json:"message"`
}

func NewError(m string) Error {
	return Error{Message: m}
}

// TODO:
type ValidationErrors struct {
	Errors string `json:"errors"`
}
