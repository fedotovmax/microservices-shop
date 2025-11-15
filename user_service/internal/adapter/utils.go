package adapter

import (
	"github.com/fedotovmax/microservices-shop/user_service/internal/ports"
)

func ParsePartial(affected, expected int64) error {
	if affected := affected; affected != expected {
		return &ports.ErrPartialUpdate{
			Expected: expected,
			Actual:   affected,
		}
	}
	return nil
}
