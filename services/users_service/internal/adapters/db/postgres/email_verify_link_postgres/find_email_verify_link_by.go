package emailverifylinkpostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/jackc/pgx/v5"
)

func findEmailVerifyLinkQuery(column db.VerifyEmailLinkEntityFields) string {
	return fmt.Sprintf("select link, user_id, link_expires_at from email_verification where %s = $1;", column)
}

func (p *postgres) FindBy(ctx context.Context, column db.VerifyEmailLinkEntityFields, value string) (*domain.EmailVerifyLink, error) {

	const op = "adapters.db.postgres.email_verify_link.FindBy"

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
			return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s:  %w: %v", op, adapters.ErrInternal, err)
	}

	return emailVerifyLink, nil
}
