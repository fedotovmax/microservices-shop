package domain

import (
	"encoding/json"
	"time"
)

type SuccessEvent struct {
	ID   string
	Type string
}

type FailedEvent struct {
	ID    string
	Type  string
	Error error
}

type CreateEvent struct {
	AggregateID string
	Topic       string
	Type        string
	Payload     json.RawMessage
}

type EventStatus string

const EventStatusNew EventStatus = "new"
const EventStatusDone EventStatus = "done"

type Event struct {
	ID          string
	AggregateID string
	Topic       string
	Type        string
	Payload     json.RawMessage
	Status      EventStatus
	CreatedAt   time.Time
	ReservedTo  *time.Time
}
