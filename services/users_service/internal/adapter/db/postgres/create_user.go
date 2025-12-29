package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
)

const createUserQuery = "insert into users (email, password_hash, created_at, updated_at) values ($1, $2, $3, $4) returning id, email;"

const createProfileQuery = "insert into profiles (user_id, updated_at) values ($1, $2);"

func (p *postgresAdapter) CreateUser(ctx context.Context, in *inputs.CreateUserInput) (*domain.UserPrimaryFields, error) {
	const op = "adapter.db.postgres.CreateUser"

	tx := p.ex.ExtractTx(ctx)

	now := time.Now()

	row := tx.QueryRow(ctx, createUserQuery, in.GetEmail(), in.GetPassword(), now, now)

	pf := &domain.UserPrimaryFields{}

	err := row.Scan(&pf.ID, &pf.Email)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	_, err = tx.Exec(ctx, createProfileQuery, pf.ID, now)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}
	return pf, nil

}
