package brandspostgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain"
)

func getAllQuery(onlyActive bool) string {

	activeFilter := ""
	if onlyActive {
		activeFilter = "and is_active = true"
	}

	return fmt.Sprintf(`
		select id, title, slug, description, logo_url, is_active, created_at, updated_at, deleted_at
		from brands where deleted_at is null
		%s
		order by title COLLATE "und-x-icu";`, activeFilter)
}

func (p *postgres) GetAll(ctx context.Context, onlyActive bool) ([]domain.Brand, error) {
	const op = "adapters.db.postgres.brands.GetAll"

	tx := p.ex.ExtractTx(ctx)

	q := getAllQuery(onlyActive)

	rows, err := tx.Query(ctx, q)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	defer rows.Close()

	var brands []domain.Brand

	for rows.Next() {

		b := domain.Brand{}

		err := rows.Scan(
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
			return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
		}

		brands = append(brands, b)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return brands, nil

}
