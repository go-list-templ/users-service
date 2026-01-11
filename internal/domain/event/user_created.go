package event

import (
	"github.com/go-list-templ/grpc/internal/domain/entity"
)

type UserCreated struct {
	Event
}

func NewUserCreated(user entity.User) UserCreated {
	return UserCreated{
		*NewEvent(
			user.ID.Value().String(),
			"user",
			nil,
		),
	}
}
