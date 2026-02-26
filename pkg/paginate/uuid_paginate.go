package paginate

import (
	"encoding/base64"

	"github.com/google/uuid"
)

type UUIDPaginate struct {
	Token string
}

func NewUUIDPaginate(token string) *UUIDPaginate {
	return &UUIDPaginate{Token: token}
}

func (u *UUIDPaginate) Limit() int {
	return DefaultLimit
}

func (u *UUIDPaginate) Cursor() string {
	if u.Token == "" {
		return ""
	}

	decodedBytes, _ := base64.URLEncoding.DecodeString(u.Token)
	decodedCursor := string(decodedBytes)
	if err := uuid.Validate(decodedCursor); err != nil {
		return ""
	}

	return decodedCursor
}

func (u *UUIDPaginate) GenerateToken(cursor string) string {
	if cursor == "" {
		return ""
	}

	return base64.URLEncoding.EncodeToString([]byte(cursor))
}
