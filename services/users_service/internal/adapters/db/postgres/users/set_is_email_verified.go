package users

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters"
)

const setIsEmailVerifiedQuery = "update users set is_email_verified = $1 where id = $1;"

func (p *postgres) SetIsEmailVerified(ctx context.Context, uid string, flag bool) error {
	const op = "adapters.db.postgres.SetIsEmailVerified"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, setIsEmailVerifiedQuery, flag, uid)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil

}
