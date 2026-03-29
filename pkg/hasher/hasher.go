package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/alexedwards/argon2id"
)

func Hash(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func Compare(hash, password string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false
	}

	return match
}

func EmailHash(email string) string {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	hash := sha256.Sum256([]byte(cleanEmail))

	return hex.EncodeToString(hash[:])
}
