package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
)

const updateSessionQuery = `update sessions set
refresh_hash = $1, ip = $2, browser = $3, browser_version = $4, 
os = $5, device = $6, expires_at = $7, updated_at = $8
where id = $9;`

func (p *postgres) UpdateSession(ctx context.Context, in *inputs.CreateSessionInput) error {

	const op = "adapter.db.postgres.UpdateSession"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, updateSessionQuery, in.RefreshHash, in.IP, in.Browser, in.BrowserVersion, in.OS, in.Device, in.ExpiresAt, time.Now(), in.SID)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil

}
