package port

import (
	"context"

	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/dto"
)

//go:generate mockgen -source=service.go -destination=mock/mock_service.go -package=mock

type (
	UserService interface {
		GetByEmail(context.Context, dto.GetByEmailInput) (entity.User, error)
		VerifyCred(context.Context, dto.VerifyCredInput) (entity.User, error)
		Create(context.Context, dto.CreateInput) (entity.User, error)
		List(context.Context, dto.ListInput) (dto.ListOutput, error)
	}
)
