package vo

import "errors"

type PasswordHash struct {
	value string
}

func NewPasswordHash(hash string) (PasswordHash, error) {
	if hash == "" {
		return PasswordHash{}, errors.New("password hash is empty")
	}
	return PasswordHash{value: hash}, nil
}

func (p PasswordHash) Value() string {
	return p.value
}
