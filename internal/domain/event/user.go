package event

import (
	"time"

	"github.com/go-list-templ/grpc/internal/domain/entity"
)

const (
	Created TypeUserEvent = "created"
	Deleted TypeUserEvent = "deleted"
)

type (
	TypeUserEvent string

	UserEvent struct {
		Event     TypeUserEvent
		Entity    entity.User
		EventTime time.Time
	}
)

func UserCreated(user entity.User) UserEvent {
	return UserEvent{
		Event:     Created,
		Entity:    user,
		EventTime: time.Now(),
	}
}
