package sessionspostgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
)

const createSessionUserQuery = `
insert into sessions_users
(uid, email)
values ($1, $2) returning uid, email;`

func (p *postgres) CreateSessionUser(ctx context.Context, uid, email string) (string, error) {

	const op = "adapter.db.postgres.CreateSessionUser"

	tx := p.ex.ExtractTx(ctx)

	row := tx.QueryRow(ctx, createSessionUserQuery, uid, email)

	var newUID string

	err := row.Scan(&newUID)

	if err != nil {
		return "", fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return newUID, nil

}
