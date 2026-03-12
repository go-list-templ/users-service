package paginate

import (
	"encoding/base64"
	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/google/uuid"
	"log"
)

type UserPaginate struct {
	token string
}

func NewUserPaginate(token string) *UserPaginate {
	return &UserPaginate{token: token}
}

func (u *UserPaginate) Token() string {
	return u.token
}

func (u *UserPaginate) Limit() int {
	return DefaultLimit + LimitOffset
}

func (u *UserPaginate) Cursor() string {
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

func (u *UserPaginate) GenerateToken(cursor string) string {
	if cursor == "" {
		return ""
	}

	return base64.URLEncoding.EncodeToString([]byte(cursor))
}

func (u *UserPaginate) NextPageToken(items []entity.User) string {
	if len(items) < u.Limit() {
		return ""
	}

	lastIndex := len(items) - (LimitOffset + 1)
	last := items[lastIndex]

	log.Println("last email: ", last.Email.Value())

	return u.GenerateToken(last.ID.Value().String())
}
