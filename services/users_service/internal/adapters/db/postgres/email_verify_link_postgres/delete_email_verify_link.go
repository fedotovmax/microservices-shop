package emailverifylinkpostgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters"
)

const deleteEmailVerifiedQuery = "delete from email_verification where link = $1;"

func (p *postgres) Delete(ctx context.Context, link string) error {

	const op = "adapters.db.postgres.email_verify_link.Delete"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, deleteEmailVerifiedQuery, link)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil

}
