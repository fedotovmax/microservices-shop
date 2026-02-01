package userspostgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
)

type buildQueryResult struct {
	Query string
	Args  []any
}

func buildUpdateUserProfileQuery(input *inputs.UpdateUserInput, id string) (*buildQueryResult, error) {

	queryParts := make([]string, 0)

	args := make([]any, 0)

	add := func(expr string, arg any) {
		queryParts = append(queryParts, fmt.Sprintf(expr, len(args)+1))
		args = append(args, arg)
	}

	if b := input.GetBirthDate(); b != nil {
		add("birth_date = $%d", b)
	}

	if f := input.GetFirstName(); f != nil {
		add("first_name = $%d", f)
	}

	if l := input.GetLastName(); l != nil {
		add("last_name = $%d", l)
	}

	if m := input.GetMiddleName(); m != nil {
		add("middle_name = $%d", m)
	}

	if url := input.GetAvatarURL(); url != nil {
		add("avatar_url = $%d", url)
	}

	if g := input.GetGender(); g != nil {
		add("gender = $%d", g)
	}

	if len(queryParts) > 0 {

		queryParts = append(queryParts, "updated_at = now()")

		query := fmt.Sprintf("update profiles set %s where user_id = $%d;", strings.Join(queryParts, ", "),
			len(args)+1)

		args = append(args, id)

		r := &buildQueryResult{
			Query: query,
			Args:  args,
		}

		return r, nil
	}
	return nil, adapters.ErrNoFieldsToUpdate
}

func (p *postgres) UpdateUserProfile(ctx context.Context, id string, in *inputs.UpdateUserInput) error {

	const op = "adapters.db.postgres.UpdateUserProfile"

	bqr, err := buildUpdateUserProfileQuery(in, id)

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
