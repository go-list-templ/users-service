package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-list-templ/grpc/pkg/pagination"

	"github.com/go-list-templ/grpc/internal/adapter/persistence/postgres"
	"github.com/go-list-templ/grpc/internal/adapter/persistence/postgres/repo/dao"
	"github.com/go-list-templ/grpc/internal/adapter/persistence/postgres/transaction"
	"github.com/go-list-templ/grpc/internal/core/domain/entity"
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

	return u.toPostgresError(err)
}

func (u *UserRepo) All(ctx context.Context, paginate pagination.Paginate) ([]entity.User, error) {
	query := `
		SELECT id, name, email, avatar, created_at, updated_at 
	   	FROM users
	   	WHERE id < $1
	   	ORDER BY id DESC 
	   	LIMIT $2
	`

	rows, err := u.Query(ctx, query, paginate.Cursor(), paginate.Limit())
	if err != nil {
		u.logger.Warn("query", zap.Error(err))

		return nil, err
	}

	users, err := pgx.CollectRows(rows, dao.RowToEntity)
	if err != nil {
		u.logger.Warn("mapping", zap.Error(err))

		return nil, err
	}

	return users, nil
}

func (u *UserRepo) toPostgresError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		u.logger.Warn("operation", zap.Error(err))

		return fmt.Errorf("operation: %w", err)
	}

	switch pgErr.Code {
	case postgres.ErrCodeAlreadyExists:
		return entity.ErrUserAlreadyExists
	case postgres.ErrCodeNotFound:
		return entity.ErrUserNotFound
	case postgres.ErrCodeInvalidData:
		return entity.ErrUserInvalidData
	default:
		infraErr := fmt.Errorf("code %s: %w", pgErr.Code, err)

		u.logger.Warn("postgres", zap.Error(infraErr))

		return infraErr
	}
}
