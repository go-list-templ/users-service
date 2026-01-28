package vo

import (
	"errors"
	"net/mail"
)

var ErrInvalidEmail = errors.New("invalid email")

type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
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
