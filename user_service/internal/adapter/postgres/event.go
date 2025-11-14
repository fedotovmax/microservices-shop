package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
	"github.com/fedotovmax/pgxtx"
)

var errEventCreate = errors.New("unexpected error when create event")

type eventPostgres struct {
	ex pgxtx.Extractor
}

func NewEventPostgres(ex pgxtx.Extractor) *eventPostgres {
	return &eventPostgres{
		ex: ex,
	}
}

func (p *eventPostgres) CreateEvent(ctx context.Context, d domain.CreateEvent) (string, error) {
	const op = "postgres.event.CreateEvent"

	tx := p.ex.ExtractTx(ctx)

	sql := "insert into events (aggregate_id, event_topic, event_type, payload) values ($1,$2,$3) returning id;"

	row := tx.QueryRow(ctx, sql, d.AggregateID, d.Topic, d.Type, d.Payload)

	var id string

	err := row.Scan(&id)

	if err != nil {
		return "", fmt.Errorf("%s: %w: %v", op, errEventCreate, err)
	}

	return id, nil

}
