package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/user_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/user_service/internal/ports"
	"github.com/fedotovmax/pgxtx"
)

type eventPostgres struct {
	ex pgxtx.Extractor
}

func NewEventPostgres(ex pgxtx.Extractor) *eventPostgres {
	return &eventPostgres{
		ex: ex,
	}
}

func (p *eventPostgres) FindNewAndNotReserved(ctx context.Context, limit int) ([]*domain.Event, error) {
	const op = "adapter.postgres.event.FindNewAndNotReserved"

	tx := p.ex.ExtractTx(ctx)

	const sql = `select id, aggregate_id, event_topic, event_type,
	payload, status, created_at, reserved_to
	from events where status != $1 AND
	(reserved_to IS NULL OR reserved_to < now());`

	rows, err := tx.Query(ctx, sql, domain.EventStatusDone)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, ports.ErrInternal, err)
	}
	defer rows.Close()

	var events []*domain.Event

	for rows.Next() {

		e := &domain.Event{}

		err := rows.Scan(&e.ID, &e.AggregateID, &e.Topic, &e.Type, &e.Payload,
			&e.Status, &e.CreatedAt, &e.ReservedTo)

		if err != nil {
			return nil, fmt.Errorf("%s: %w: %v", op, ports.ErrInternal, err)
		}

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, ports.ErrInternal, err)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ports.ErrNotFound)
	}

	return events, nil

}

func (p *eventPostgres) Create(ctx context.Context, d domain.CreateEvent) (string, error) {
	const op = "adapter.postgres.event.CreateEvent"

	tx := p.ex.ExtractTx(ctx)

	const sql = `insert into events (aggregate_id, event_topic, event_type, payload)
	values ($1,$2,$3,$4) returning id;`

	row := tx.QueryRow(ctx, sql, d.AggregateID, d.Topic, d.Type, d.Payload)

	var id string

	err := row.Scan(&id)

	if err != nil {
		return "", fmt.Errorf("%s: %w: %v", op, ports.ErrInternal, err)
	}

	return id, nil
}

func (p *eventPostgres) RemoveReserve(ctx context.Context, id string) error {

	const op = "adapter.postgres.event.ChangeStatus"

	const sql = "update events set status = $1 where id = $2;"

	tx := p.ex.ExtractTx(ctx)

	result, err := tx.Exec(ctx, sql, domain.EventStatusDone, id)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, ports.ErrInternal, err)
	}

	var expected int64 = 1

	if err := adapter.ParsePartial(result.RowsAffected(), expected); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *eventPostgres) ChangeStatus(ctx context.Context, id string) error {
	const op = "adapter.postgres.event.ChangeStatus"

	const sql = "update events set status = $1 where id = $2;"

	tx := p.ex.ExtractTx(ctx)

	result, err := tx.Exec(ctx, sql, domain.EventStatusDone, id)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, ports.ErrInternal, err)
	}

	var expected int64 = 1

	if err := adapter.ParsePartial(result.RowsAffected(), expected); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *eventPostgres) SetReservedToByIDs(ctx context.Context, ids []string, dur time.Duration) error {

	const op = "adapter.postgres.event.SetReservedToByIds"

	reservedTo := time.Now().Add(dur)

	const sql = "update events set reserved_to = $1 where id = ANY ($2);"

	tx := p.ex.ExtractTx(ctx)

	res, err := tx.Exec(ctx, sql, reservedTo, ids)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, ports.ErrInternal, err)
	}

	expected := int64(len(ids))

	if affected := res.RowsAffected(); affected != expected {
		return fmt.Errorf("%s: %w", op, &ports.ErrPartialUpdate{
			Expected: expected,
			Actual:   affected,
		})
	}

	return nil

}
