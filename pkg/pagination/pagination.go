package pagination

import (
	"encoding/base64"
	"errors"
)

const (
	DefaultLimitList = 10
	MaxLimitList     = 100
)

var (
	ErrEmptyCursor   = errors.New("empty cursor")
	ErrInvalidCursor = errors.New("invalid cursor")
)

type Paginate struct {
	PageSize int64
	Token    string
}

func New(pageSize int64, token string) *Paginate {
	return &Paginate{PageSize: pageSize, Token: token}
}

func (p *Paginate) Limit() int {
	limit := int(p.PageSize)

	if limit <= 0 {
		limit = DefaultLimitList
	}

	if limit >= MaxLimitList {
		limit = MaxLimitList
	}

	return limit
}

func (p *Paginate) Cursor() string {
	if p.Token == "" {
		return ""
	}

	return base64.StdEncoding.EncodeToString([]byte(p.Token))
}

func (p *Paginate) GenerateToken(cursor string) string {
	if cursor == "" {
		return ""
	}

	token, _ := base64.StdEncoding.DecodeString(cursor)
	return string(token)
}
