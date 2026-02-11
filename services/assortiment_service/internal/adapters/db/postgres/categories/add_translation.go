package categories

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
)

const createTranslationQuery = "insert into category_translations (category_id, language_code, title, description) values ($1, $2, $3, $4);"

func (p *postgres) AddTranslation(
	ctx context.Context,
	categoryID string,
	in *inputs.AddCategoryTranslate,
) error {
	const op = "adapters.db.postgres.categories.AddTranslation"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(
		ctx,
		createTranslationQuery,
		categoryID,
		in.LanguageCode,
		in.Title,
		in.Description,
	)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil

}
