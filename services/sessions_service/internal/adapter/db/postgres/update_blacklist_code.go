package postgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
)

const updateBlacklistCodeQuery = `
update blacklist (code, code_expires_at) values ($1, $2) where uid = $3;`

func (p *postgres) UpdateBlacklistCode(ctx context.Context, in *inputs.SecurityInput) error {

	const op = "adapter.db.postgres.UpdateBlacklistCode"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, updateBlacklistCodeQuery, in.Code, in.CodeExpiresAt, in.UID)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil

}
