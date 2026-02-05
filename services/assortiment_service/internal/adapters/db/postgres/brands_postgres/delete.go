package brandspostgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters/db"
)

func deleteQuery(column db.BrandEntityFields) string {
	return fmt.Sprintf("update brands set deleted_at = $1 where %s = $2;", column)
}

func (p *postgres) Delete(ctx context.Context, column db.BrandEntityFields, searchValue string) error {

	const op = "adapters.db.postgres.brands.Delete"

	err := db.IsBrandEntityField(column)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tx := p.ex.ExtractTx(ctx)

	_, err = tx.Exec(ctx, deleteQuery(column), time.Now().UTC(), searchValue)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil

}
