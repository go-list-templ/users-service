package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func Compare(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func EmailHash(email string) string {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	hash := sha256.Sum256([]byte(cleanEmail))

	return hex.EncodeToString(hash[:])
}
