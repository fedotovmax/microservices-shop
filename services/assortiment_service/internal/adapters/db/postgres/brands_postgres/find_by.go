package brandspostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain"
	"github.com/jackc/pgx/v5"
)

func findByQuery(column db.BrandEntityFields, onlyActive bool) string {
	activeFilter := ""
	if onlyActive {
		activeFilter = "and is_active = true"
	}
	return fmt.Sprintf(`
	select id, title, slug, description, logo_url, is_active,
	created_at, updated_at, deleted_at from brands
	where %s = $1
	%s
	;
	`, column, activeFilter)
}

func (p *postgres) FindBy(ctx context.Context, params *db.FindBrandByParams) (
	*domain.Brand, error,
) {
	const op = "adapters.db.postgres.brands.FindBy"

	err := db.IsBrandEntityField(params.SearchColumn)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tx := p.ex.ExtractTx(ctx)

	q := findByQuery(params.SearchColumn, params.OnlyActive)

	row := tx.QueryRow(ctx, q, params.SearchValue)

	b := &domain.Brand{}

	err = row.Scan(
		&b.ID,
		&b.Title,
		&b.Slug,
		&b.Description,
		&b.LogoURL,
		&b.IsActive,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return b, nil

}
