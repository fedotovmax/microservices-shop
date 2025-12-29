package postgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
)

const removeUserFromBlacklistQuery = "delete from blacklist where uid = $1;"

func (p *postgres) RemoveUserFromBlacklist(ctx context.Context, uid string) error {

	const op = "adapter.db.postgres.RemoveUserFromBlacklist"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, removeUserFromBlacklistQuery, uid)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil
}
