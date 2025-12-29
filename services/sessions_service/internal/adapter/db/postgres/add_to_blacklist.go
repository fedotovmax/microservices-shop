package postgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
)

const addToBlackListQuery = `
insert into blacklist (uid, code, code_expires_at) values ($1, $2, $3);`

func (p *postgres) AddToBlackList(ctx context.Context, in *inputs.AddToBlackListInput) error {

	const op = "adapter.db.postgres.AddToBlacklist"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, addToBlackListQuery, in.UID, in.Code, in.CodeExpiresAt)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil

}
