package dao

import (
	"time"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/domain/vo"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/mo"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      *string   `db:"name"`
	Password  string    `db:"password"`
	Email     string    `db:"email"`
	Avatar    string    `db:"avatar"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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

func RowToEntity(row pgx.CollectableRow) (entity.User, error) {
	d, err := pgx.RowToStructByNameLax[User](row)
	if err != nil {
		return entity.User{}, err
	}

	return d.ToEntity(), nil
}
