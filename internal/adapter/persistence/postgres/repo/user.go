package repo

import (
	"context"

	"github.com/go-list-templ/grpc/internal/adapter/persistence/postgres"
	"github.com/go-list-templ/grpc/internal/adapter/persistence/postgres/repo/dao"
	"github.com/go-list-templ/grpc/internal/adapter/persistence/postgres/transaction"
	"github.com/go-list-templ/grpc/internal/core/domain/entity"
)

type UserRepo struct {
	*postgres.Postgres

	getter *transaction.TrmGetter
}

func NewUserRepo(p *postgres.Postgres, g *transaction.TrmGetter) *UserRepo {
	return &UserRepo{p, g}
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

	return err
}

func (u *UserRepo) All(ctx context.Context) ([]entity.User, error) {
	var users []entity.User

	rows, err := u.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var userDAO dao.User

		err = rows.Scan(&userDAO)
		if err != nil {
			return nil, err
		}

		users = append(users, userDAO.ToEntity())
	}

	return users, nil
}
