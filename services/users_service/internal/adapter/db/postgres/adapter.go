package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/pgxtx"
	"github.com/jackc/pgx/v5"
)

type postgresAdapter struct {
	ex  pgxtx.Extractor
	log *slog.Logger
}

func NewPostgresAdapter(ex pgxtx.Extractor, log *slog.Logger) *postgresAdapter {

	return &postgresAdapter{
		ex:  ex,
		log: log,
	}
}

func (p *postgresAdapter) FindUserBy(ctx context.Context, column db.UserEntityFields, value string) (*domain.User, error) {

	const op = "adapter.postgres.FindUserBy"

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
			return nil, fmt.Errorf("%s: %s: %w: %v", op, column, adapter.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %s: %w: %v", op, column, adapter.ErrInternal, err)
	}

	return u, nil
}

func (p *postgresAdapter) UpdateUserProfile(ctx context.Context, id string, in *inputs.UpdateUserInput) error {

	const op = "adapter.postgres.UpdateUserProfile"

	bqr, err := buildUpdateUserProfileQuery(in, id)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tx := p.ex.ExtractTx(ctx)

	_, err = tx.Exec(ctx, bqr.Query, bqr.Args...)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil

}

func (p *postgresAdapter) CreateUser(ctx context.Context, in *inputs.CreateUserInput) (*domain.UserPrimaryFields, error) {
	const op = "adapter.postgres.CreateUser"

	tx := p.ex.ExtractTx(ctx)

	row := tx.QueryRow(ctx, createUserQuery, in.GetEmail(), in.GetPassword())

	pf := &domain.UserPrimaryFields{}

	err := row.Scan(&pf.ID, &pf.Email)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	_, err = tx.Exec(ctx, createProfileQuery, pf.ID)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}
	return pf, nil

}

func (p *postgresAdapter) CreateEmailVerifyLink(ctx context.Context, userID string) (*domain.EmailVerifyLink, error) {

	const op = "adapter.postgres.CreateEmailVerifyLink"

	tx := p.ex.ExtractTx(ctx)

	emailVerifyLink := &domain.EmailVerifyLink{}

	row := tx.QueryRow(ctx, createEmailVerifyLinkQuery, userID)

	err := row.Scan(&emailVerifyLink.Link, &emailVerifyLink.UserID, &emailVerifyLink.ValidityPeriod)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return emailVerifyLink, nil

}

func (p *postgresAdapter) FindEmailVerifyLink(ctx context.Context, link string) (*domain.EmailVerifyLink, error) {

	const op = "adapter.postgres.FindEmailVerifyLink"

	tx := p.ex.ExtractTx(ctx)

	emailVerifyLink := &domain.EmailVerifyLink{}

	row := tx.QueryRow(ctx, findEmailVerifyLinkQuery, link)

	err := row.Scan(&emailVerifyLink.Link, &emailVerifyLink.UserID, &emailVerifyLink.ValidityPeriod)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s:  %w: %v", op, adapter.ErrInternal, err)
	}

	return emailVerifyLink, nil
}

func (p *postgresAdapter) UpdateEmailVerifyLinkByUserID(ctx context.Context, userID string) (*domain.EmailVerifyLink, error) {

	const op = "adapter.postgres.UpdateEmailVerifyLinkByUserID"

	tx := p.ex.ExtractTx(ctx)

	emailVerifyLink := &domain.EmailVerifyLink{}

	row := tx.QueryRow(ctx, updateEmailVerifyLinkByUserIDQuery, userID)

	err := row.Scan(&emailVerifyLink.Link, &emailVerifyLink.UserID, &emailVerifyLink.ValidityPeriod)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return emailVerifyLink, nil

}

//TODO:
// func (p *postgresAdapter) FindEvents(ctx context.Context, f any) ([]*domain.Event, error) {
// 	const op = "adapter.postgres.FindEvents"

// 	return nil, nil
// }

func (p *postgresAdapter) FindNewAndNotReservedEvents(ctx context.Context, limit int) ([]*domain.Event, error) {

	const op = "adapter.postgres.FindNewAndNotReservedEvents"

	tx := p.ex.ExtractTx(ctx)

	rows, err := tx.Query(ctx, findNewAndNotReservedEventsQuery, domain.EventStatusDone, limit)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}
	defer rows.Close()

	var events []*domain.Event

	for rows.Next() {

		e := &domain.Event{}

		err := rows.Scan(&e.ID, &e.AggregateID, &e.Topic, &e.Type, &e.Payload,
			&e.Status, &e.CreatedAt, &e.ReservedTo)

		if err != nil {
			return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
		}

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return events, nil

}

func (p *postgresAdapter) CreateEvent(ctx context.Context, in *inputs.CreateEvent) (string, error) {
	const op = "outbox.postgres.CreateEvent"

	tx := p.ex.ExtractTx(ctx)

	row := tx.QueryRow(ctx, createEventQuery,
		in.GetAggregateID(), in.GetTopic(), in.GetType(), in.GetPayload())

	var id string

	err := row.Scan(&id)

	if err != nil {
		return "", fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return id, nil
}

func (p *postgresAdapter) RemoveEventReserve(ctx context.Context, id string) error {

	const op = "adapter.postgres.RemoveEventReserve"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, removeEventReserveQuery, id)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil
}

func (p *postgresAdapter) SetEventStatusDone(ctx context.Context, id string) error {
	const op = "adapter.postgres.SetEventStatusDone"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, setEventStatusDoneQuery, domain.EventStatusDone, id)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil
}

func (p *postgresAdapter) SetEventsReservedToByIDs(ctx context.Context, ids []string, dur time.Duration) error {

	const op = "adapter.postgres.SetEventsReservedToByIDs"

	reservedTo := time.Now().Add(dur)

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, setEventsReservedToByIDsQuery, reservedTo, ids)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil

}
