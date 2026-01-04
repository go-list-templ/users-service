package storage

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/repo/storage/dao"
	"github.com/go-list-templ/grpc/pkg/postgres"
)

type UserPostgresRepo struct {
	*postgres.Postgres
}

func NewUserPostgresRepo(postgres *postgres.Postgres) *UserPostgresRepo {
	return &UserPostgresRepo{postgres}
}

func (r *UserPostgresRepo) Store(ctx context.Context, user entity.User) error {
	query := `
		INSERT INTO users (id, name, email, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.Exec(ctx, query,
		user.ID.Value(),
		user.Name.Value(),
		user.Email.Value(),
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserPostgresRepo) All(ctx context.Context) ([]entity.User, error) {
	var users []entity.User

	rows, err := r.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var userDAO dao.User

		err = rows.Scan(&userDAO.ID, &userDAO.Name, &userDAO.Email, &userDAO.CreatedAt, &userDAO.UpdatedAt)
		if err != nil {
			return nil, err
		}

		user, err := userDAO.ToEntity()
		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return users, nil
}
