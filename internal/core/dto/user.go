package dto

import (
	"time"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
)

type (
	User struct {
		ID        string
		Name      *string
		Email     string
		Avatar    string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	VerifyCredInput struct {
		Email    string
		Password string
	}

	GetByEmailInput struct {
		Email string
	}

	CreateInput struct {
		Name     *string
		Email    string
		Password string
	}

	ListInput struct {
		PageToken string
	}

	ListOutput struct {
		Users         []User
		NextPageToken string
	}
)

func FromEntity(user entity.User) User {
	var username *string
	if name, ok := user.Name.Get(); ok {
		str := name.Value()
		username = &str
	}

	return User{
		ID:        user.ID.Value().String(),
		Name:      username,
		Email:     user.Email.Value(),
		Avatar:    user.Avatar.Value(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func FromEntities(entities []entity.User) []User {
	users := make([]User, len(entities))

	for i, user := range entities {
		users[i] = FromEntity(user)
	}

	return users
}
