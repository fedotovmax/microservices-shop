package postgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
	"github.com/fedotovmax/pgxtx"
)

type userPostgres struct {
	ex pgxtx.Extractor
}

func NewUserPostgres(ex pgxtx.Extractor) *userPostgres {
	return &userPostgres{
		ex: ex,
	}
}

func (p *userPostgres) GetByID(ctx context.Context, id string) error {
	return nil
}

func (p *userPostgres) CreateUser(ctx context.Context, d domain.CreateUser) (string, error) {
	const op = "adapter.postgres.user.CreateUser"

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
