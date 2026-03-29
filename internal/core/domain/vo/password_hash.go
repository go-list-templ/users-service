package vo

import "github.com/go-list-templ/users-service/pkg/hasher"

type PasswordHash struct {
	value string
}

func NewPasswordHash(password PlainPassword) (PasswordHash, error) {
	hash, err := hasher.Hash(password.Value())
	if err != nil {
		return PasswordHash{}, err
	}

	return PasswordHash{value: hash}, nil
}

func (p PasswordHash) Value() string {
	return p.value
}

func (p PasswordHash) Compare(password PlainPassword) bool {
	return hasher.Compare(p.value, password.Value())
}

func UnsafePasswordHash(hash string) PasswordHash {
	return PasswordHash{value: hash}
}
