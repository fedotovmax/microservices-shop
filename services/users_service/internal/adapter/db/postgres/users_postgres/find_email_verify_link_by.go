package userspostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/jackc/pgx/v5"
)

// const findEmailVerifyLinkQuery = "select link, user_id, link_expires_at from email_verification where link = $1;"

func findEmailVerifyLinkQuery(column db.VerifyEmailLinkEntityFields) string {
	return fmt.Sprintf("select link, user_id, link_expires_at from email_verification where %s = $1;", column)
}

func (p *postgres) FindEmailVerifyLinkBy(ctx context.Context, column db.VerifyEmailLinkEntityFields, value string) (*domain.EmailVerifyLink, error) {

	const op = "adapter.db.postgres.FindEmailVerifyLink"

	err := db.IsVerifyEmailEntityField(column)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tx := p.ex.ExtractTx(ctx)

	emailVerifyLink := &domain.EmailVerifyLink{}

	row := tx.QueryRow(ctx, findEmailVerifyLinkQuery(column), value)

	err = row.Scan(&emailVerifyLink.Link, &emailVerifyLink.UserID, &emailVerifyLink.LinkExpiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s:  %w: %v", op, adapter.ErrInternal, err)
	}

	return emailVerifyLink, nil
}
