package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/jackc/pgx/v5"
)

const findUserQuery = "select su.uid, su.email, b.code, b.code_expires_at from sessions_users as su left join blacklist as b on su.uid = b.uid where su.uid = $1;"

func (p *postgres) FindUser(ctx context.Context, uid string) (*domain.SessionsUser, error) {
	const op = "adapter.db.postgres.FindUser"

	tx := p.ex.ExtractTx(ctx)

	row := tx.QueryRow(ctx, findUserQuery, uid)

	var code *string
	var expiresAt *time.Time
	user := &domain.SessionsUser{}

	err := row.Scan(&user.Info.UID, &user.Info.Email, &code, &expiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	if code != nil && expiresAt != nil {
		user.BlackList = &domain.BlackList{
			Code:          *code,
			CodeExpiresAt: *expiresAt,
		}
	}

	return user, nil
}
