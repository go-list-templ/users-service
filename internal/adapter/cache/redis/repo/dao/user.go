package dao

import (
	"time"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/domain/vo"
	"github.com/google/uuid"
	"github.com/samber/mo"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      *string   `json:"name,omitempty"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromEntity(user entity.User) User {
	var username *string
	if name, ok := user.Name.Get(); ok {
		str := name.Value()
		username = &str
	}

	return User{
		ID:        user.ID.Value(),
		Name:      username,
		Password:  user.Password.Value(),
		Email:     user.Email.Value(),
		Avatar:    user.Email.Value(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *User) ToEntity() entity.User {
	username := mo.None[vo.Name]()
	if name, ok := mo.PointerToOption(u.Name).Get(); ok {
		username = mo.Some(vo.UnsafeName(name))
	}

	return entity.User{
		ID:        vo.UnsafeID(u.ID),
		Name:      username,
		Password:  vo.UnsafePasswordHash(u.Password),
		Email:     vo.UnsafeEmail(u.Email),
		Avatar:    vo.UnsafeAvatar(u.Avatar),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) IsEmpty() bool {
	return u.ID == uuid.Nil
}
