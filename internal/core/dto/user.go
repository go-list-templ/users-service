package dto

import (
	"time"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
)

type (
	User struct {
		ID        string
		Name      string
		Email     string
		Avatar    string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	UserCreateInput struct {
		Name  string
		Email string
	}

	UserCreateOutput struct {
		User User
	}

	UserListInput struct {
		PageToken string
	}

	UserListOutput struct {
		Users         []User
		NextPageToken string
	}
)

func UserFromEntity(user entity.User) User {
	return User{
		ID:        user.ID.Value().String(),
		Name:      user.Name.Value(),
		Email:     user.Email.Value(),
		Avatar:    user.Avatar.Value(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UsersFromEntity(entities []entity.User) []User {
	users := make([]User, len(entities))

	for i, user := range entities {
		users[i] = UserFromEntity(user)
	}

	return users
}
