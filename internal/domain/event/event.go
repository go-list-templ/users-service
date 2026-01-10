package event

import "time"

type Event interface {
	Payload() []byte
	OccurredAt() time.Time
	AggregateID() string
	AggregateType() string
}
