package userspostgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
)

const deleteEmailVerifiedQuery = "delete from email_verification where link = $1;"

func (p *postgres) DeleteEmailVerifyLink(ctx context.Context, link string) error {

	const op = "adapter.db.postgres.DeleteEmailVerifyLink"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, deleteEmailVerifiedQuery, link)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil

}
