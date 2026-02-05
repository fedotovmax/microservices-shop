package brandspostgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters/db"
)

type buildUpdateQueryResult struct {
	Query string
	Args  []any
}

func updateQuery(params *db.UpdateBrandParams) (*buildUpdateQueryResult, error) {

	err := db.IsBrandEntityField(params.SearchColumn)

	if err != nil {
		return nil, err
	}

	queryParts := make([]string, 0)

	args := make([]any, 0)

	add := func(expr string, arg any) {
		queryParts = append(queryParts, fmt.Sprintf(expr, len(args)+1))
		args = append(args, arg)
	}

	if params.Input.Title != nil {
		add("title = $%d", *params.Input.Title)
	}

	if params.NewSlug != nil {
		add("slug = $%d", *params.NewSlug)
	}

	if params.Input.Description != nil {
		add("description = $%d", *params.Input.Description)
	}

	if params.Input.LogoURL != nil {
		add("logo_url = $%d", *params.Input.LogoURL)
	}

	if params.Input.IsActive != nil {
		add("is_active = $%d", *params.Input.IsActive)
	}

	if len(queryParts) > 0 {

		add("updated_at = $%d", time.Now().UTC())

		q := fmt.Sprintf(
			"update brands set %s where %s = $%d;",
			strings.Join(queryParts, ", "),
			params.SearchColumn,
			len(args)+1,
		)

		args = append(args, params.SearchValue)

		r := &buildUpdateQueryResult{
			Query: q,
			Args:  args,
		}

		return r, nil

	}

	return nil, adapters.ErrNoFieldsToUpdate
}

func (p *postgres) Update(
	ctx context.Context,
	params *db.UpdateBrandParams,
) error {

	const op = "adapters.db.postgres.brands.Update"

	bqr, err := updateQuery(params)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tx := p.ex.ExtractTx(ctx)

	_, err = tx.Exec(ctx, bqr.Query, bqr.Args...)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}
