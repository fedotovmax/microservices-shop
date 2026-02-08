package badgespostgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	postgresPkg "github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters/db/postgres"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
)

func updateQuery(in *inputs.UpdateBadge) (*postgresPkg.BuildUpdateQueryResult, error) {

	queryParts := make([]string, 0)

	args := make([]any, 0)

	add := func(expr string, arg any) {
		queryParts = append(queryParts, fmt.Sprintf(expr, len(args)+1))
		args = append(args, arg)
	}

	if in.Color != nil {
		add("color = $%d", *in.Color)
	}

	if in.Priority != nil {
		add("priority = $%d", *in.Priority)
	}

	if in.EndsAt != nil {
		add("ends_at = $%d", *in.EndsAt)
	}

	if len(queryParts) > 0 {
		q := fmt.Sprintf(
			"update badges set %s where code = $%d;",
			strings.Join(queryParts, ", "),
			len(args)+1,
		)

		args = append(args, in.Code)

		r := &postgresPkg.BuildUpdateQueryResult{
			Query: q,
			Args:  args,
		}

		return r, nil

	}

	return nil, adapters.ErrNoFieldsToUpdate

}

func (p *postgres) Update(ctx context.Context, in *inputs.UpdateBadge) error {

	const op = "adapters.db.postgres.badges.Update"

	tx := p.ex.ExtractTx(ctx)

	br, err := updateQuery(in)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, br.Query, br.Args...)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}
