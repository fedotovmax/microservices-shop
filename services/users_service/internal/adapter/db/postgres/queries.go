package postgres

import (
	"fmt"
	"strings"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
)

//USER QUERIES

func findUserByQuery(column db.UserEntityFields) string {
	return fmt.Sprintf(`select u.id, u.email, u.phone, u.password_hash, u.created_at, u.updated_at,
p.last_name, p.first_name, p.middle_name, p.birth_date, p.gender, p.avatar_url, p.updated_at
from users as u
inner join profiles as p on u.id = p.user_id
where u.deleted_at is null and u.%s = $1;
`, column)
}

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
		add("gender = $%d", g.String())
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
	return nil, adapter.ErrNoFieldsToUpdate
}

const createUserQuery = "insert into users (email, password_hash) values ($1, $2) returning id;"

const createProfileQuery = "insert into profiles (user_id) values ($1);"

//EVENT QUERIES

const findNewAndNotReservedEventsQuery = `select id, aggregate_id, event_topic, event_type,
	payload, status, created_at, reserved_to
	from events where status != $1 AND
	(reserved_to IS NULL OR reserved_to < now())
	order by created_at asc
	limit $2;`

const setEventsReservedToByIDsQuery = "update events set reserved_to = $1 where id = ANY ($2);"

const removeEventReserveQuery = "update events set reserved_to = null where id = $1;"

const createEventQuery = `insert into events (aggregate_id, event_topic, event_type, payload)
values ($1,$2,$3,$4) returning id;`

const setEventStatusDoneQuery = "update events set status = $1 where id = $2;"
