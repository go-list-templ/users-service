package vo

import (
	"fmt"
	"strings"
)

const (
	MinLengthPass = 8
	MaxLengthPass = 30
)

var (
	ErrPlainPasswordMinLength = fmt.Errorf("min pass must be at least %v characters", MinLengthPass)
	ErrPlainPasswordMaxLength = fmt.Errorf("max pass must be at least %v characters", MaxLengthPass)
)

type PlainPassword struct {
	value string
}

func NewPlainPassword(password string) (PlainPassword, error) {
	password = strings.TrimSpace(password)

	if len(password) < MinLengthPass {
		return PlainPassword{}, ErrPlainPasswordMinLength
	}
	if len(password) > MaxLengthPass {
		return PlainPassword{}, ErrPlainPasswordMaxLength
	}

	return PlainPassword{value: password}, nil
}

func (u *PlainPassword) Value() string {
	return u.value
}
