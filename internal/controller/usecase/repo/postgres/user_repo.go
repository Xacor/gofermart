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
	const (
		sqlInsertUser = `INSERT INTO users(
			login, password) 
			VALUES($1,$2);`

		sqlInsertBalance = `INSERT INTO balances(
			current, withdrawn, user_id)
			VALUES ($1, $2, $3);`

		sqlQueryUser = "SELECT id FROM users WHERE login = $1;"
	)

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("can not begin transaction: %v", err)
	}

	_, err = tx.Exec(ctx, sqlInsertUser, user.Login, user.Password)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("can not create user error: %v", err)
	}

	var userID int
	err = tx.QueryRow(ctx, sqlQueryUser, user.Login).Scan(&userID)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("can not query user_id: %v", err)
	}

	_, err = tx.Exec(ctx, sqlInsertBalance, 0, 0, userID)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("can not create balance: %v", err)
	}

	return tx.Commit(ctx)
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
	const sql = "SELECT * FROM users WHERE login = $1;"

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("cannot GetByLogin error: %v", err)
	}

	return user, nil
}
