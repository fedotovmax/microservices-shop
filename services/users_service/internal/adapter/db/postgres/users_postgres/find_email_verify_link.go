package userspostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/jackc/pgx/v5"
)

const findEmailVerifyLinkQuery = "select link, user_id, link_expires_at from email_verification where link = $1;"

func (p *postgres) FindEmailVerifyLink(ctx context.Context, link string) (*domain.EmailVerifyLink, error) {

	const op = "adapter.db.postgres.FindEmailVerifyLink"

	tx := p.ex.ExtractTx(ctx)

	emailVerifyLink := &domain.EmailVerifyLink{}

	row := tx.QueryRow(ctx, findEmailVerifyLinkQuery, link)

	err := row.Scan(&emailVerifyLink.Link, &emailVerifyLink.UserID, &emailVerifyLink.LinkExpiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s:  %w: %v", op, adapter.ErrInternal, err)
	}

	return emailVerifyLink, nil
}
