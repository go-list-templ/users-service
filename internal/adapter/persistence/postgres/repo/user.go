package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-list-templ/users-service/internal/adapter/persistence/postgres"
	"github.com/go-list-templ/users-service/internal/adapter/persistence/postgres/repo/dao"
	"github.com/go-list-templ/users-service/internal/adapter/persistence/postgres/transaction"
	"github.com/go-list-templ/users-service/internal/core/domain/entity"
	"github.com/go-list-templ/users-service/internal/core/domain/entityerr"
	"github.com/go-list-templ/users-service/pkg/paginate"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type UserRepo struct {
	*postgres.Postgres

	logger *zap.Logger
	getter *transaction.TrmGetter
}

func NewUserRepo(p *postgres.Postgres, l *zap.Logger, g *transaction.TrmGetter) *UserRepo {
	return &UserRepo{p, l, g}
}

func (u *UserRepo) Store(ctx context.Context, user entity.User) error {
	query := `
		INSERT INTO users (id, name, email, avatar, created_at, updated_at) 
       	VALUES (@id, @name, @email, @avatar, @created_at, @updated_at)
	`

	args := pgx.NamedArgs{
		"id":         user.ID.Value(),
		"name":       user.Name.Value(),
		"email":      user.Email.Value(),
		"avatar":     user.Avatar.Value(),
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	_, err := u.getter.TrOrDB(ctx, u.Postgres).Exec(ctx, query, args)

	return u.toPostgresError(ctx, err)
}

func (u *UserRepo) All(ctx context.Context, paginate paginate.Paginate) ([]entity.User, error) {
	var args []any

	query := `
		SELECT id, name, email, avatar, created_at, updated_at 
	   	FROM users
	   	WHERE id < $1
	   	ORDER BY id DESC 
	   	LIMIT $2
	`

	cursor := paginate.Cursor()
	limit := paginate.Limit()

	args = []any{cursor, limit}

	if len(cursor) == 0 {
		query = `
		SELECT id, name, email, avatar, created_at, updated_at 
	   	FROM users
	   	ORDER BY id DESC 
	   	LIMIT $1
		`

		args = []any{limit}
	}

	rows, err := u.Query(ctx, query, args...)
	if err != nil {
		return nil, u.toPostgresError(ctx, err)
	}

	users, err := pgx.CollectRows(rows, dao.RowToEntity)
	if err != nil {
		return nil, u.toPostgresError(ctx, err)
	}

	return users, nil
}

func (u *UserRepo) toPostgresError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		u.logger.Error("operation", zap.Any("context", ctx), zap.Error(err))

		return fmt.Errorf("operation: %w", err)
	}

	switch pgErr.Code {
	case postgres.ErrCodeAlreadyExists:
		return entityerr.ErrUserAlreadyExists
	case postgres.ErrCodeNotFound:
		return entityerr.ErrUserNotFound
	case postgres.ErrCodeInvalidData:
		return entityerr.ErrUserInvalidData
	default:
		infraErr := fmt.Errorf("code %s: %w", pgErr.Code, err)

		u.logger.Error("postgres", zap.Any("context", ctx), zap.Error(infraErr))

		return infraErr
	}
}
