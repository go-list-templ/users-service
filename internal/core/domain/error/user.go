package error

import (
	"errors"
	"fmt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserInvalidData   = errors.New("invalid user data")
)

type UserError struct {
	Field string
	Err   error
}

func NewUserError(field string, err error) *UserError {
	return &UserError{Field: field, Err: err}
}

func (u *UserError) Error() string {
	return fmt.Sprintf("invalid user %s: %v", u.Field, u.Err)
}

func (u *UserError) Unwrap() error {
	return ErrUserInvalidData
}
