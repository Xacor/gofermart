package postgres

import (
	"context"
	"fmt"

	"github.com/Xacor/gophermart/internal/entity"
	"github.com/Xacor/gophermart/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) CreateUser(ctx context.Context, user entity.User) error {
	const sql = "INSERT INTO users (login, password) VALUES($1,$2);"
	_, err := r.Pool.Exec(ctx, sql, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("cannot create user error: %v", err)
	}

	return nil
}
func (r *UserRepo) GetByID(ctx context.Context, id int) (entity.User, error) {
	const sql = "SELECT * FROM users WHERE id = $1;"

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, id).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
func (r *UserRepo) GetByLogin(ctx context.Context, login string) (entity.User, error) {
	const sql = "SELECT * FROM users WHERE users.login = $1;"

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("cannot GetByLogin error: %v", err)
	}

	return user, nil
}
