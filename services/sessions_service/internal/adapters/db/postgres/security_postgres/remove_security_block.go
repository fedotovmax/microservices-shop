package securitypostgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db"
)

const removeUserFromBlacklistQuery = "delete from blacklist where uid = $1;"

const removeUserBypassQuery = "delete from bypass where uid = $1;"

func (p *postgres) RemoveSecurityBlock(ctx context.Context, table db.SecurityTable, uid string) error {

	const op = "adapter.db.postgres.RemoveSecurityBlock"

	var err error

	err = db.IsSecurityTable(table)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tx := p.ex.ExtractTx(ctx)

	switch table {
	case db.SecurityTableBlacklist:
		_, err = tx.Exec(ctx, removeUserFromBlacklistQuery, uid)
	case db.SecurityTableBypass:
		_, err = tx.Exec(ctx, removeUserBypassQuery, uid)
	}

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}
