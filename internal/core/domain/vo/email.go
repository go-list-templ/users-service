package vo

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidEmail = errors.New("invalid email")

	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	if !emailRegex.MatchString(email) {
		return Email{}, ErrInvalidEmail
	}

	return Email{value: email}, nil
}

func UnsafeEmail(email string) Email {
	return Email{value: email}
}

func (e *Email) Value() string {
	return e.value
}
