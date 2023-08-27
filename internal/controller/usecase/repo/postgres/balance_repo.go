package postgres

import (
	"context"

	"github.com/Xacor/gophermart/internal/entity"
	"github.com/Xacor/gophermart/pkg/postgres"
)

type BalanceRepo struct {
	*postgres.Postgres
}

// AddBonuses implements usecase.BalanceRepo.
func (*BalanceRepo) AddBonuses(ctx context.Context, userID int, amount int) error {
	panic("unimplemented")
}

// Get implements usecase.BalanceRepo.
func (*BalanceRepo) Get(ctx context.Context, userID int) (entity.Balance, error) {
	panic("unimplemented")
}

// Withdraw implements usecase.BalanceRepo.
func (*BalanceRepo) Withdraw(ctx context.Context, userID int, amount int) error {
	panic("unimplemented")
}

func NewBalanceRepo(pg *postgres.Postgres) *BalanceRepo {
	return &BalanceRepo{pg}
}
