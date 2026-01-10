package event

import (
	"time"

	"github.com/google/uuid"
)

type UserCreated struct {
	eventID       uuid.UUID
	payload       []byte
	occurredAt    time.Time
	aggregateID   string
	aggregateType string
}

func NewUserCreated(aggregateID, aggregateType string, payload []byte) UserCreated {
	return UserCreated{
		eventID:       uuid.New(),
		payload:       payload,
		occurredAt:    time.Now(),
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
	}
}

func (e *UserCreated) Payload() []byte {
	return e.payload
}

func (e *UserCreated) OccurredAt() time.Time {
	return e.occurredAt
}

func (e *UserCreated) AggregateID() string {
	return e.aggregateID
}

func (e *UserCreated) AggregateType() string {
	return e.aggregateType
}
