package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/jackc/pgx/v5"
)

func findSessionQuery(column db.SessionEntityFields) string {
	return fmt.Sprintf(`
	select s.id, s.refresh_hash, s.ip, s.browser,
	s.browser_version, s.os, s.device, s.created_at,
	s.revoked_at, s.expires_at, s.updated_at,
	u.uid, u.email,
	b.code, b.code_expires_at 
	from sessions as s
	inner join sessions_users as u on u.uid = s.uid
	left join blacklist as b on b.uid = u.uid
	where %s = $1
	`, column)
}

func (p *postgres) FindSession(ctx context.Context, column db.SessionEntityFields, value string) (*domain.Session, error) {

	const op = "adapter.db.postgres.FindSession"

	err := db.IsSessionEntityField(column)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tx := p.ex.ExtractTx(ctx)

	row := tx.QueryRow(ctx, findSessionQuery(column), value)

	s := &domain.Session{}

	var blacklistCode *string
	var blacklistCodeExpiresAt *time.Time

	err = row.Scan(
		&s.ID,
		&s.RefreshHash,
		&s.IP,
		&s.Browser,
		&s.BrowserVersion,
		&s.OS,
		&s.Device,
		&s.CreatedAt,
		&s.RevokedAt,
		&s.ExpiresAt,
		&s.UpdatedAt,
		&s.User.Info.UID,
		&s.User.Info.Email,
		&blacklistCode,
		&blacklistCodeExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	if blacklistCode != nil && blacklistCodeExpiresAt != nil {
		s.User.BlackList = &domain.BlackList{
			Code:          *blacklistCode,
			CodeExpiresAt: *blacklistCodeExpiresAt,
		}
	}

	return s, nil
}
