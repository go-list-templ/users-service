package dao

import (
	"time"

	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/go-list-templ/grpc/internal/core/domain/vo"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromEntity(e entity.User) User {
	return User{
		e.ID.Value(),
		e.Name.Value(),
		e.Email.Value(),
		e.Avatar.Value(),
		e.CreatedAt,
		e.UpdatedAt,
	}
}

func (u *User) ToEntity() entity.User {
	return entity.User{
		ID:        vo.UnsafeID(u.ID),
		Name:      vo.UnsafeName(u.Name),
		Email:     vo.UnsafeEmail(u.Email),
		Avatar:    vo.UnsafeAvatar(u.Avatar),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
