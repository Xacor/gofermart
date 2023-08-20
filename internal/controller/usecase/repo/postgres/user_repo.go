package postgres

import (
	"context"

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
	return nil
}
func (r *UserRepo) GetByID(ctx context.Context, id string) (entity.User, error) {
	return entity.User{}, nil
}
func (r *UserRepo) GetByLogin(ctx context.Context, login string) (entity.User, error) {
	return entity.User{}, nil
}
