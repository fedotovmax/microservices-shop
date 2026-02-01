package usecases

import (
	"context"
	"fmt"
)

func (u *usecases) RevokeSessions(ctx context.Context, sids []string) error {

	const op = "usecases.security.RevokeSessions"

	err := u.sessionsStorage.RevokeSessions(ctx, sids)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
