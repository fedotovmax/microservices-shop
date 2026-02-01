package sessionspostgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters"
)

const createUserQuery = `
insert into sessions_users (uid, email) values ($1, $2);`

func (p *postgres) CreateUser(ctx context.Context, uid string, email string) error {

	const op = "adapter.db.postgres.CreateUser"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, createUserQuery, uid, email)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil

}
