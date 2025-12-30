package postgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
)

const setUserIPBypassQuery = "insert into bypass (uid, code, bypass_expires_at) values ($1, $2, $3);"

func (p *postgres) AddIPBypass(ctx context.Context, in *inputs.SecurityInput) error {

	const op = "adapter.db.postgres.AddIPBypass"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, setUserIPBypassQuery, in.UID, in.Code, in.CodeExpiresAt)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil

}
