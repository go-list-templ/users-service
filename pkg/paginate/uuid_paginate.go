package paginate

import (
	"encoding/base64"

	"github.com/google/uuid"
)

type UUIDPaginate struct {
	token string
}

func NewUUIDPaginate(token string) *UUIDPaginate {
	return &UUIDPaginate{token: token}
}

func (u *UUIDPaginate) Token() string {
	return u.token
}

func (u *UUIDPaginate) Limit() int {
	return DefaultLimit + LimitOffset
}

func (u *UUIDPaginate) Cursor() string {
	if u.token == "" {
		return ""
	}

	decodedBytes, err := base64.URLEncoding.DecodeString(u.token)
	if err != nil {
		return ""
	}

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
