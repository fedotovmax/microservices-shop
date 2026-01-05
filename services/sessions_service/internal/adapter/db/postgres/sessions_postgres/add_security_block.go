package sessionspostgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
)

const addToBlackListQuery = `
insert into blacklist (uid, code, code_expires_at) values ($1, $2, $3);`

const updateBlacklistCodeQuery = `
update blacklist 
set code = $1, code_expires_at = $2 where uid = $3;`

const addUserIPBypassQuery = "insert into bypass (uid, code, bypass_expires_at) values ($1, $2, $3);"

const updateIPBypassQuery = `
update bypass 
set code = $1, code_expires_at = $2 where uid = $3;`

func (p *postgres) AddSecurityBlock(ctx context.Context, operation db.Operation, table db.SecurityTable, in *inputs.SecurityInput) error {

	const op = "adapter.db.postgres.AddSecurityBlock"

	var err error

	err = db.IsSecurityTable(table)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = db.IsOperation(operation)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tx := p.ex.ExtractTx(ctx)

	switch table {
	case db.SecurityTableBlacklist:
		switch operation {
		case db.OperationInsert:
			_, err = tx.Exec(ctx, addToBlackListQuery, in.UID, in.Code, in.CodeExpiresAt)
		case db.OperationUpdate:
			_, err = tx.Exec(ctx, updateBlacklistCodeQuery, in.Code, in.CodeExpiresAt, in.UID)
		default:
			return fmt.Errorf("%s: %w", op, adapter.ErrUnsupported)
		}
	case db.SecurityTableBypass:
		switch operation {
		case db.OperationInsert:
			_, err = tx.Exec(ctx, addUserIPBypassQuery, in.UID, in.Code, in.CodeExpiresAt)
		case db.OperationUpdate:
			_, err = tx.Exec(ctx, updateIPBypassQuery, in.Code, in.CodeExpiresAt, in.UID)
		default:
			return fmt.Errorf("%s: %w", op, adapter.ErrUnsupported)
		}
	}

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
