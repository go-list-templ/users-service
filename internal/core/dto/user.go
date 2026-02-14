package dto

import (
	"github.com/google/uuid"
	"time"
)

type (
	User struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Avatar    string    `json:"avatar"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	UserCreateInput struct {
		Name  string
		Email string
	}

	UserListInput struct {
		PageSize  int64
		PageToken string
	}
)
