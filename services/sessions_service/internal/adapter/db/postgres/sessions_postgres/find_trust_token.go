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

const findTrustTokenQuery = "select token_hash, uid, last_used_at, expires_at, revoked_at from trust_tokens where token_hash = $1 and uid = $2 and revoked_at is null and expires_at > $3;"

func (p *postgres) FindTrustToken(ctx context.Context, uid, tokenHash string) (*domain.DeviceTrustToken, error) {
	const op = "adapter.db.postgres.FindTrustToken"

	tx := p.ex.ExtractTx(ctx)

	row := tx.QueryRow(ctx, findTrustTokenQuery, tokenHash, uid, time.Now().UTC())

	trustToken := &domain.DeviceTrustToken{}

	err := row.Scan(
		&trustToken.TokenHash, &trustToken.UID,
		&trustToken.LastUsed, &trustToken.ExpiresAt,
		&trustToken.RevokedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return trustToken, nil

}
