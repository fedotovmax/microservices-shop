package userspostgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
)

const updateEmailVerifyLinkByUserIDQuery = "update email_verification set link = gen_random_uuid(), link_expires_at = $1 where user_id = $2 returning link, user_id, link_expires_at;"

func (p *postgres) UpdateEmailVerifyLinkByUserID(ctx context.Context, uid string, expiresAt time.Time) (*domain.EmailVerifyLink, error) {

	const op = "adapter.db.postgres.UpdateEmailVerifyLinkByUserID"

	tx := p.ex.ExtractTx(ctx)

	emailVerifyLink := &domain.EmailVerifyLink{}

	row := tx.QueryRow(ctx, updateEmailVerifyLinkByUserIDQuery, expiresAt, uid)

	err := row.Scan(&emailVerifyLink.Link, &emailVerifyLink.UserID, &emailVerifyLink.LinkExpiresAt)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return emailVerifyLink, nil

}
