package domain

import (
	"encoding/json"
	"time"
)

type EventStatus string

func (es EventStatus) String() string {
	return string(es)
}

const EventStatusNew EventStatus = "new"
const EventStatusDone EventStatus = "done"

type Event struct {
	ReservedTo  *time.Time
	Payload     json.RawMessage
	CreatedAt   time.Time
	ID          string
	AggregateID string
	Topic       string
	Type        string
	Status      EventStatus
}

func (e *Event) GetID() string {
	return e.ID
}

func (e *Event) GetAggregateID() string {
	return e.AggregateID
}
func (e *Event) GetTopic() string {
	return e.Topic
}
func (e *Event) GetType() string {
	return e.Type
}
func (e *Event) GetPayload() json.RawMessage {
	return e.Payload
}
