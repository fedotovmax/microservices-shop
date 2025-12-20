package domain

import "github.com/fedotovmax/microservices-shop/notify_service/internal/domain/errs"

type Cmd string

func (c Cmd) String() string {
	return string(c)
}

func (c Cmd) Validate() error {
	switch c {
	case Start, Help:
		return nil
	default:
		return errs.ErrInvalidCommand
	}
}

const (
	Start Cmd = "/start"
	Help  Cmd = "/help"
)
