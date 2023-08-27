package postgres

import (
	"context"
	"fmt"

	"github.com/Xacor/gophermart/internal/entity"
	"github.com/Xacor/gophermart/pkg/postgres"
)

type BalanceRepo struct {
	*postgres.Postgres
}

func NewBalanceRepo(pg *postgres.Postgres) *BalanceRepo {
	return &BalanceRepo{pg}
}

// Get implements usecase.BalanceRepo.
func (r *BalanceRepo) Get(ctx context.Context, userID int) (entity.Balance, error) {
	const sql = "SELECT user_id, current, withdrawn FROM balances WHERE user_id = $1;"

	var balance entity.Balance
	err := r.Pool.QueryRow(ctx, sql, userID).Scan(&balance.UserID, &balance.Current, &balance.Withdrawn)
	if err != nil {
		return entity.Balance{}, fmt.Errorf("can not query balance: %v", err)
	}
	return balance, nil
}
