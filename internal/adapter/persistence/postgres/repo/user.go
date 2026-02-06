package repo

import (
	"context"
	"errors"
	"fmt"

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
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := u.getter.TrOrDB(ctx, u.Postgres).
		Exec(ctx, query,
			user.ID.Value(),
			user.Name.Value(),
			user.Email.Value(),
			user.Avatar.Value(),
			user.CreatedAt,
			user.UpdatedAt,
		)

	return u.toPostgresError(err)
}

func (u *UserRepo) All(ctx context.Context) ([]entity.User, error) {
	rows, err := u.Query(ctx, "SELECT * FROM users")
	if err != nil {
		u.logger.Warn("query", zap.Error(err))

		return []entity.User{}, err
	}

	users, err := pgx.CollectRows(rows, dao.RowToEntity)
	if err != nil {
		u.logger.Warn("mapping", zap.Error(err))

		return []entity.User{}, err
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
