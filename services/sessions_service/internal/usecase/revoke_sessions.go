package usecase

import (
	"context"
	"fmt"
)

func (u *usecases) RevokeSessions(ctx context.Context, sids []string) error {

	const op = "usecases.RevokeSessions"

	err := u.storage.RevokeSessions(ctx, sids)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
