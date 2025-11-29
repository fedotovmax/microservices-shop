package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/pgxtx"
	"github.com/jackc/pgx/v5"
)

type postgresAdapter struct {
	ex pgxtx.Extractor
}

func NewPostgresAdapter(ex pgxtx.Extractor) *postgresAdapter {
	return &postgresAdapter{
		ex: ex,
	}
}

func (p *postgresAdapter) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return p.findBy(ctx, "id", id)
}

func (p *postgresAdapter) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return p.findBy(ctx, "email", email)
}

func (p *postgresAdapter) findBy(ctx context.Context, column string, value string) (*domain.User, error) {

	const op = "adapter.postgres.user.findBy"

	tx := p.ex.ExtractTx(ctx)

	row := tx.QueryRow(ctx, findByQuery(column), value)

	u := &domain.User{}

	err := row.Scan(&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt, &u.Profile.LastName, &u.Profile.FirstName, &u.Profile.MiddleName, &u.Profile.BirthDate, &u.Profile.Gender, &u.Profile.AvatarURL, &u.Profile.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %s: %w: %v", op, column, adapter.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %s: %w: %v", op, column, adapter.ErrInternal, err)
	}

	return u, nil
}

func (p *postgresAdapter) Create(ctx context.Context, d domain.CreateUserInput) (string, error) {
	const op = "adapter.postgres.user.Create"

	tx := p.ex.ExtractTx(ctx)

	row := tx.QueryRow(ctx, createUserQuery, d.GetEmail(), d.GetPassword())

	var id string

	err := row.Scan(&id)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, createProfileQuery, id)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}
