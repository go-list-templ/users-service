package storage

import (
	"context"

	"github.com/go-list-templ/grpc/internal/domain/entity"
	"github.com/go-list-templ/grpc/internal/repo/storage/dao"
	"github.com/go-list-templ/grpc/pkg/postgres"
)

type UserPostgres struct {
	*postgres.Postgres
}

func NewUserPostgres(postgres *postgres.Postgres) *UserPostgres {
	return &UserPostgres{postgres}
}

func (u *UserPostgres) Store(ctx context.Context, user entity.User) error {
	query := `
		INSERT INTO users (id, name, email, avatar, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := u.Exec(ctx, query,
		user.ID.Value(),
		user.Name.Value(),
		user.Email.Value(),
		user.Avatar.Value(),
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserPostgres) All(ctx context.Context) ([]entity.User, error) {
	var users []entity.User

	rows, err := u.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var userDAO dao.User

		err = rows.Scan(
			&userDAO.ID,
			&userDAO.Name,
			&userDAO.Email,
			&userDAO.Avatar,
			&userDAO.CreatedAt,
			&userDAO.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, userDAO.ToEntity())
	}

	return users, nil
}
