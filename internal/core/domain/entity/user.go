package entity

import (
	"time"

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

func NewUser(Name, email string) (*User, error) {
	validName, err := vo.NewName(Name)
	if err != nil {
		return nil, err
	}

	validEmail, err := vo.NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        vo.NewID(),
		Name:      validName,
		Email:     validEmail,
		Avatar:    vo.NewAvatar(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}
