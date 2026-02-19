package paginate

import (
	"encoding/base64"
	"github.com/google/uuid"
)

type UUIDPaginate struct {
	PageSize int64
	Token    string
}

func NewUUIDPaginate(pageSize int64, token string) *UUIDPaginate {
	return &UUIDPaginate{PageSize: pageSize, Token: token}
}

func (u *UUIDPaginate) Limit() int {
	limit := int(u.PageSize)

	if limit <= 0 {
		limit = DefaultLimitList
	}

	if limit >= MaxLimitList {
		limit = MaxLimitList
	}

	return limit
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
