package userspostgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
)

const createEmailVerifyLinkQuery = "insert into email_verification (user_id, link_expires_at) values ($1, $2) returning link, user_id, link_expires_at;"

func (p *postgres) CreateEmailVerifyLink(ctx context.Context, userID string, expiresAt time.Time) (*domain.EmailVerifyLink, error) {

	const op = "adapter.db.postgres.CreateEmailVerifyLink"

	tx := p.ex.ExtractTx(ctx)

	emailVerifyLink := &domain.EmailVerifyLink{}

	row := tx.QueryRow(ctx, createEmailVerifyLinkQuery, userID, expiresAt)

	err := row.Scan(&emailVerifyLink.Link, &emailVerifyLink.UserID, &emailVerifyLink.LinkExpiresAt)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return emailVerifyLink, nil

}
