package dao

import (
	"time"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/domain/vo"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FromEntity(e entity.User) User {
	return User{
		e.ID.Value(),
		e.Name.Value(),
		e.Email.Value(),
		e.CreatedAt,
		e.UpdatedAt,
	}
}

func (u *User) ToEntity() entity.User {
	return entity.User{
		ID:        vo.UnsafeID(u.ID),
		Name:      vo.UnsafeName(u.Name),
		Email:     vo.UnsafeEmail(u.Email),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
