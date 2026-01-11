package event

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID            uuid.UUID
	Payload       []byte
	OccurredAt    time.Time
	AggregateID   string
	AggregateType string
}

func NewEvent(aggregateID, aggregateType string, payload []byte) *Event {
	return &Event{
		ID:            uuid.New(),
		Payload:       payload,
		OccurredAt:    time.Now(),
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
	}
}
