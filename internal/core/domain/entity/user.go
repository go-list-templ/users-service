package entity

import (
	"time"

	entityErr "github.com/go-list-templ/grpc/internal/core/domain/error"

	"github.com/go-list-templ/grpc/internal/core/domain/vo"
)

type User struct {
	ID        vo.ID
	Name      vo.Name
	Email     vo.Email
	Avatar    vo.Avatar
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name, email string) (User, error) {
	id, err := vo.NewID()
	if err != nil {
		return User{}, err
	}

	validName, err := vo.NewName(name)
	if err != nil {
		return User{}, entityErr.NewUserError("name", err)
	}

	validEmail, err := vo.NewEmail(email)
	if err != nil {
		return User{}, entityErr.NewUserError("email", err)
	}

	return User{
		ID:        id,
		Name:      validName,
		Email:     validEmail,
		Avatar:    vo.NewAvatar(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}
