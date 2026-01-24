package event

import (
	"encoding/json"
	"github.com/go-list-templ/grpc/internal/domain/entity"
)

type UserCreated struct {
	Event
}

func NewUserCreated(user entity.User) (UserCreated, error) {
	payload, err := json.Marshal(user)
	if err != nil {
		return UserCreated{}, err
	}

	return UserCreated{
		*NewEvent(
			user.ID.Value().String(),
			"user_created",
			payload,
		),
	}, nil
}
