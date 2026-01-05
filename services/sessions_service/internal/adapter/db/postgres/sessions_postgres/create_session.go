package sessionspostgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
)

const createSessionQuery = `
insert into sessions
(id, uid, refresh_hash, ip, browser, browser_version, os, device, expires_at, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id;`

func (p *postgres) CreateSession(ctx context.Context, in *inputs.CreateSessionInput) (string, error) {

	const op = "adapter.db.postgres.CreateSession"

	tx := p.ex.ExtractTx(ctx)

	now := time.Now()

	row := tx.QueryRow(ctx, createSessionQuery, in.SID, in.UID, in.RefreshHash, in.IP, in.Browser, in.BrowserVersion, in.OS, in.Device, in.ExpiresAt, now, now)

	var sessionId string

	err := row.Scan(&sessionId)

	if err != nil {
		return "", fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return sessionId, nil

}
