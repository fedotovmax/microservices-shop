package brands

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
)

const createQuery = "insert into brands (title, slug, description, logo_url, created_at, updated_at) values ($1, $2, $3, $4, $5, $6);"

func (p *postgres) Create(ctx context.Context, in *inputs.CreateBrand, slug string) error {
	const op = "adapters.db.postgres.brands.Create"

	tx := p.ex.ExtractTx(ctx)

	now := time.Now().UTC()

	_, err := tx.Exec(ctx, createQuery, in.Title, slug, in.Description, in.LogoURL, now, now)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil

}
