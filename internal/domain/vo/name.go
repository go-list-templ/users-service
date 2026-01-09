package vo

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	MinLength = 5
	MaxLength = 30
)

var (
	ErrNameMinLength = fmt.Errorf("name must be at least %v characters", MinLength)
	ErrNameMaxLength = fmt.Errorf("name must be at least %v characters", MinLength)
	ErrNameValidate  = fmt.Errorf("name can only contain letters, numbers and underscores")

	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

type Name struct {
	value string
}

func NewName(name string) (Name, error) {
	name = strings.TrimSpace(name)

	if len(name) < MinLength {
		return Name{}, ErrNameMinLength
	}
	if len(name) > MaxLength {
		return Name{}, ErrNameMaxLength
	}

	if !nameRegex.MatchString(name) {
		return Name{}, ErrNameValidate
	}

	return Name{value: name}, nil
}

func UnsafeName(name string) Name {
	return Name{value: name}
}

func (u *Name) Value() string {
	return u.value
}
