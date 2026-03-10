package event

import (
	"encoding/json"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/domain/event/dao"
)

const AggregateUserCreated = "user_created"

type UserCreated struct {
	Event
}

func NewUserCreated(user entity.User) (UserCreated, error) {
	userDAO := dao.User{
		ID:        user.ID.Value(),
		Name:      user.Name.Value(),
		Email:     user.Email.Value(),
		Avatar:    user.Avatar.Value(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	payload, err := json.Marshal(userDAO)
	if err != nil {
		return UserCreated{}, err
	}

	return UserCreated{
		*NewEvent(
			user.ID.Value().String(),
			AggregateUserCreated,
			payload,
		),
	}, nil
}
