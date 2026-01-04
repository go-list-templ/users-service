package dao

import (
	"time"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/domain/vo"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *User) ToEntity() (*entity.User, error) {
	id, err := vo.NewIDFromString(u.ID.String())
	if err != nil {
		return nil, err
	}

	name, err := vo.NewName(u.Name)
	if err != nil {
		return nil, err
	}

	email, err := vo.NewEmail(u.Email)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:        id,
		Name:      name,
		Email:     email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}
