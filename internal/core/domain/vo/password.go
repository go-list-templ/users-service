package vo

import (
	"fmt"
	"strings"
)

const (
	MinLengthPass = 2
	MaxLengthPass = 30
)

var (
	ErrPasswordMinLength = fmt.Errorf("min pass must be at least %v characters", MinLengthPass)
	ErrPasswordMaxLength = fmt.Errorf("max pass must be at least %v characters", MaxLengthPass)
)

type Password struct {
	value string
}

func NewPassword(password string) (Password, error) {
	password = strings.TrimSpace(password)

	if len(password) < MinLengthPass {
		return Password{}, ErrPasswordMinLength
	}
	if len(password) > MaxLengthPass {
		return Password{}, ErrPasswordMaxLength
	}

	//todo add hash for pass

	return Password{value: password}, nil
}

func (u *Password) Value() string {
	return u.value
}
