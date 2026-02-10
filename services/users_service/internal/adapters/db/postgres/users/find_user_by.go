package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/jackc/pgx/v5"
)

//u.deleted_at is null and

func findUserByQuery(column db.UserEntityFields) string {
	return fmt.Sprintf(`select u.id, u.email, u.phone, u.password_hash, u.is_email_verified,
u.is_phone_verified, u.created_at, u.updated_at, u.deleted_at,
p.last_name, p.first_name, p.middle_name, p.birth_date, p.gender, p.avatar_url, p.updated_at
from users as u
inner join profiles as p on u.id = p.user_id
where u.%s = $1;
`, column)
}

func (p *postgres) FindBy(ctx context.Context, column db.UserEntityFields, value string) (*domain.User, error) {

	const op = "adapters.db.postgres.FindBy"

	err := db.IsUserEntityField(column)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tx := p.ex.ExtractTx(ctx)

	row := tx.QueryRow(ctx, findUserByQuery(column), value)

	u := &domain.User{}

	err = row.Scan(&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.IsEmailVerified, &u.IsPhoneVerified, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt, &u.Profile.LastName, &u.Profile.FirstName, &u.Profile.MiddleName, &u.Profile.BirthDate, &u.Profile.Gender, &u.Profile.AvatarURL, &u.Profile.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %s: %w: %v", op, column, adapters.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %s: %w: %v", op, column, adapters.ErrInternal, err)
	}

	return u, nil
}
