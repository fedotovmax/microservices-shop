package categoriespostgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
)

const createQuery = "insert into categories (parent_id, slug, logo_url, created_at, updated_at) values ($1,$2,$3,$4,$5) returning id;"

func (p *postgres) Create(
	ctx context.Context,
	slug string,
	parentID *string,
	logoURL *string,
) (string, error) {

	const op = "adapters.db.postgres.categories.Create"

	tx := p.ex.ExtractTx(ctx)

	now := time.Now().UTC()

	row := tx.QueryRow(ctx, createQuery, parentID, slug, logoURL, now, now)

	var categoryID string

	err := row.Scan(&categoryID)

	if err != nil {
		return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return categoryID, nil

}
