package badgespostgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
)

const createTranslationQuery = "insert into badge_translations (badge_code, language_code, title) values ($1, $2, $3);"

func (p *postgres) AddTranslation(
	ctx context.Context,
	badgeCode string,
	in *inputs.BadgeTranslate,
) error {
	const op = "adapters.db.postgres.badges.AddTranslation"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(
		ctx,
		createTranslationQuery,
		badgeCode,
		in.LanguageCode,
		in.Title,
	)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil

}
