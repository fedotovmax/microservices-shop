package sessions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/jackc/pgx/v5"
)

// TODO: add migration and scan deleted_at field
func findSessionQuery(column db.SessionEntityFields) string {
	return fmt.Sprintf(`
	select s.id, s.refresh_hash, s.ip, s.browser,
	s.browser_version, s.os, s.device, s.created_at,
	s.revoked_at, s.expires_at, s.updated_at,
	u.uid, u.email,
	bl.code as bl_code, bl.code_expires_at as bl_code_expires_at,
	bp.code as bp_code, bp.bypass_expires_at as bp_code_expires_at
	from sessions as s
	inner join sessions_users as su on su.uid = s.uid
	left join blacklist as bl on bl.uid = su.uid
	left join bypass as bp on su.uid = bp.uid
	where %s = $1
	`, column)
}

func (p *postgres) FindBy(ctx context.Context, column db.SessionEntityFields, value string) (*domain.Session, error) {

	const op = "adapter.db.postgres.sessions.FindBy"

	err := db.IsSessionEntityField(column)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tx := p.ex.ExtractTx(ctx)

	row := tx.QueryRow(ctx, findSessionQuery(column), value)

	s := &domain.Session{}
	user := &domain.SessionsUser{}

	var blacklistCode *string
	var blacklistCodeExpiresAt *time.Time

	var bypassCode *string
	var bypassCodeExpiresAt *time.Time

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
		&user.Info.UID,
		&user.Info.Email,
		&blacklistCode,
		&blacklistCodeExpiresAt,
		&bypassCode,
		&bypassCodeExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	if blacklistCode != nil && blacklistCodeExpiresAt != nil {
		user.BlackList = &domain.BlackList{
			Code:          *blacklistCode,
			CodeExpiresAt: *blacklistCodeExpiresAt,
		}
	}

	if bypassCode != nil && bypassCodeExpiresAt != nil {
		user.Bypass = &domain.Bypass{
			Code:            *bypassCode,
			BypassExpiresAt: *bypassCodeExpiresAt,
		}
	}

	s.User = user

	return s, nil
}
