package postgres

import "fmt"

func findByQuery(column string) string {
	return fmt.Sprintf(`select u.id, u.email, u.phone, u.password_hash, u.created_at, u.updated_at,
p.last_name, p.first_name, p.middle_name, p.birth_date, p.gender, p.avatar_url, p.updated_at
from users as u
inner join profiles as p on u.id = p.user_id
where u.deleted_at is null and u.%s = $1;
`, column)
}

const createUserQuery = "insert into users (email, password_hash) values ($1, $2) returning id;"

const createProfileQuery = "insert into profiles (user_id) values ($1);"
