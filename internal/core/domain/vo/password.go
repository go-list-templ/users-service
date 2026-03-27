package vo

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	MinLengthPass = 2
	MaxLengthPass = 30
)

var (
	ErrPasswordMinLength = fmt.Errorf("min pass must be at least %v characters", MinLengthPass)
	ErrPasswordMaxLength = fmt.Errorf("max pass must be at least %v characters", MaxLengthPass)
	ErrPasswordValidate  = fmt.Errorf("pass can only contain letters, numbers and underscores")

	passRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
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

	if !passRegex.MatchString(password) {
		return Password{}, ErrPasswordValidate
	}

	return Password{value: password}, nil
}

func UnsafePassword(password string) Password {
	return Password{value: password}
}

func (u *Password) Value() string {
	return u.value
}
