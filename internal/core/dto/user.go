package dto

import "github.com/go-list-templ/users-service/internal/core/domain/entity"

type (
	UserCreateInput struct {
		Name  string
		Email string
	}

	UserListInput struct {
		PageToken string
	}

	UserListOutput struct {
		Users         []entity.User
		NextPageToken string
	}
)
