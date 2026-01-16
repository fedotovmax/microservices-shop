package eventspostgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
)

const findNewAndNotReservedEventsQuery = `select id, aggregate_id, event_topic, event_type,
	payload, status, created_at, reserved_to
	from events where status != $1 AND
	(reserved_to IS NULL OR reserved_to < $2)
	order by created_at asc
	limit $3;`

func (p *postgres) FindNewAndNotReservedEvents(ctx context.Context, limit int) ([]*domain.Event, error) {

	const op = "adapter.db.postgres.FindNewAndNotReservedEvents"

	tx := p.ex.ExtractTx(ctx)

	rows, err := tx.Query(ctx, findNewAndNotReservedEventsQuery, domain.EventStatusDone,
		time.Now().UTC(), limit)

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
