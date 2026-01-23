package sessionspostgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/jackc/pgx/v5"
)

// TODO: add migration and scan deleted_at field
const findUserQuery = `
select su.uid, su.email, bl.code, bl.code_expires_at, bp.code, bp.bypass_expires_at
from sessions_users as su
left join blacklist as bl on su.uid = bl.uid
left join bypass as bp on su.uid = bp.uid
where su.uid = $1;`

func (p *postgres) FindUser(ctx context.Context, uid string) (*domain.SessionsUser, error) {

	const op = "adapter.db.postgres.FindUser"

	tx := p.ex.ExtractTx(ctx)

	row := tx.QueryRow(ctx, findUserQuery, uid)

	var blCode *string
	var blExpiresAt *time.Time

	var bpCode *string
	var bpExpiresAt *time.Time

	user := &domain.SessionsUser{}

	err := row.Scan(&user.Info.UID, &user.Info.Email, &blCode, &blExpiresAt, &bpCode, &bpExpiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	if blCode != nil && blExpiresAt != nil {
		user.BlackList = &domain.BlackList{
			Code:          *blCode,
			CodeExpiresAt: *blExpiresAt,
		}
	}

	if bpCode != nil && bpExpiresAt != nil {
		user.Bypass = &domain.Bypass{
			Code:            *bpCode,
			BypassExpiresAt: *bpExpiresAt,
		}
	}

	return user, nil
}
