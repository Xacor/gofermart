package postgres

import (
	"context"
	"time"

	"github.com/Xacor/gophermart/internal/entity"
	"github.com/Xacor/gophermart/pkg/postgres"
)

type WithdrawalsRepo struct {
	*postgres.Postgres
}

func NewWithdrawalsRepo(pg *postgres.Postgres) *WithdrawalsRepo {
	return &WithdrawalsRepo{pg}
}

// Create implements usecase.WithdrawalsRepo.
func (r *WithdrawalsRepo) Create(ctx context.Context, withdraw entity.Withdraw) error {
	const (
		sqlInsertWithdraw = "INSERT INTO withdrawals(order_id, user_id, sum, processed_at) VALUES ($1, $2, $3, $4);"
		sqlBalance        = "UPDATE public.balances SET current = current - $1, withdrawn = withdrawn + $1 WHERE user_id = $2;"
	)
	withdraw.ProcessedAt = time.Now()
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sqlInsertWithdraw, withdraw.Order, withdraw.UserID, withdraw.Sum, withdraw.ProcessedAt)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, sqlBalance, withdraw.Sum, withdraw.UserID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

// GetByUserID implements usecase.WithdrawalsRepo.
func (r *WithdrawalsRepo) GetByUserID(ctx context.Context, userID int) ([]entity.Withdraw, error) {
	const sql = "SELECT id, user_id, order_id, sum, processed_at FROM withdrawals WHERE user_id = $1;"
	var withdrawals []entity.Withdraw
	rows, err := r.Pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var withdraw entity.Withdraw
		err = rows.Scan(&withdraw.ID, &withdraw.UserID, &withdraw.Order, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdraw)
	}

	return withdrawals, nil
}
