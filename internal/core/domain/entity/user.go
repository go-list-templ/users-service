package entity

import (
	"time"

	"github.com/go-list-templ/users-service/internal/core/domain/entityerr"
	"github.com/go-list-templ/users-service/internal/core/domain/vo"
	"github.com/samber/mo"
)

type User struct {
	ID        vo.ID
	Name      mo.Option[vo.Name]
	Password  vo.Password
	Email     vo.Email
	Avatar    vo.Avatar
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name *string, email, password string) (User, error) {
	id, err := vo.NewID()
	if err != nil {
		return User{}, err
	}

	validEmail, err := vo.NewEmail(email)
	if err != nil {
		return User{}, entityerr.NewUserError("email", err)
	}

	validPass, err := vo.NewPassword(password)
	if err != nil {
		return User{}, entityerr.NewUserError("password", err)
	}

	validName := mo.None[vo.Name]()
	if name, ok := mo.PointerToOption(name).Get(); ok {
		v, err := vo.NewName(name)
		if err != nil {
			return User{}, entityerr.NewUserError("name", err)
		}

		validName = mo.Some(v)
	}

	return User{
		ID:        id,
		Name:      validName,
		Email:     validEmail,
		Password:  validPass,
		Avatar:    vo.NewAvatar(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}
