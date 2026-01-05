package sessionspostgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
)

const findUserSessionsQuery = `
  select s.id, s.refresh_hash, s.ip, s.browser,
	s.browser_version, s.os, s.device, s.created_at,
	s.revoked_at, s.expires_at, s.updated_at,
	su.uid, su.email,
	bl.code, bl.code_expires_at,
	bp.code, bp.bypass_expires_at
	from sessions as s
	inner join sessions_users as su on su.uid = s.uid
	left join blacklist as bl on bl.uid = su.uid
	left join bypass as bp on su.uid = bp.uid
  where su.uid = $1 order by s.updated_at desc;
`

func (p *postgres) FindUserSessions(ctx context.Context, uid string) ([]*domain.Session, error) {

	const op = "adapter.db.postgres.FindUserSessions"

	tx := p.ex.ExtractTx(ctx)

	rows, err := tx.Query(ctx, findUserSessionsQuery, uid)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	defer rows.Close()

	var sessions []*domain.Session

	for rows.Next() {

		s := &domain.Session{}

		var blacklistCode *string
		var blacklistCodeExpiresAt *time.Time

		var bypassCode *string
		var bypassCodeExpiresAt *time.Time

		err = rows.Scan(
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
			&bypassCode,
			&bypassCodeExpiresAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
		}

		if blacklistCode != nil && blacklistCodeExpiresAt != nil {
			s.User.BlackList = &domain.BlackList{
				Code:          *blacklistCode,
				CodeExpiresAt: *blacklistCodeExpiresAt,
			}
		}

		if bypassCode != nil && bypassCodeExpiresAt != nil {
			s.User.Bypass = &domain.Bypass{
				Code:            *bypassCode,
				BypassExpiresAt: *bypassCodeExpiresAt,
			}
		}

		sessions = append(sessions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return sessions, nil

}
