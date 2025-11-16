package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/user_service/internal/ports"
	"github.com/fedotovmax/pgxtx"
	"github.com/jackc/pgx/v5"
)

type userPostgres struct {
	ex pgxtx.Extractor
}

func NewUserPostgres(ex pgxtx.Extractor) *userPostgres {
	return &userPostgres{
		ex: ex,
	}
}

func (p *userPostgres) FindByID(ctx context.Context, id string) (*domain.User, error) {

	const op = "adapter.postgres.user.FindByID"

	tx := p.ex.ExtractTx(ctx)

	const sql = "select id, email, first_name, last_name from users where id = $1;"

	row := tx.QueryRow(ctx, sql, id)

	u := &domain.User{}

	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w: %v", op, ports.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w: %v", op, ports.ErrInternal, err)
	}

	return u, nil
}

func (p *userPostgres) Create(ctx context.Context, d domain.CreateUser) (string, error) {
	const op = "adapter.postgres.user.Create"

	const sql = "insert into users (email, first_name, last_name) values ($1, $2, $3) returning id;"

	tx := p.ex.ExtractTx(ctx)

	row := tx.QueryRow(ctx, sql, d.Email, d.FirstName, d.LastName)

	var id string

	err := row.Scan(&id)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}
