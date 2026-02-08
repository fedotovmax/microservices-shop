package categoriespostgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
)

const updateTranslationQuery = `
update category_translations
set title = $1, description = $2
where id = $3;`

func (p *postgres) UpdateTranslation(ctx context.Context, in *inputs.UpdateCategoryTranslate) error {
	const op = "adapters.db.postgres.categories.UpdateTranslation"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, updateTranslationQuery, in.Title, in.Description, in.ID)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil

}
