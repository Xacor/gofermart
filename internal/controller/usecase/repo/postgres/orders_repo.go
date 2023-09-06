package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Xacor/gophermart/internal/entity"
	"github.com/Xacor/gophermart/pkg/postgres"
)

var (
	ErrNoRows = errors.New("no rows")
)

type OrderRepo struct {
	*postgres.Postgres
}

func NewOrderRepo(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{pg}
}

// Create implements usecase.OrderRepo.
func (r *OrderRepo) Create(ctx context.Context, order entity.Order) error {
	const sql = `INSERT INTO public.orders(
	id, user_id, status, accrual, uploaded_at)
	VALUES ($1, $2, $3, $4, $5);`

	order.UploadedAt = time.Now()

	_, err := r.Pool.Exec(ctx, sql, order.Number, order.UserID, order.Status, order.Accrual, order.UploadedAt)
	if err != nil {
		return fmt.Errorf("cannot create order error: %v", err)
	}

	return nil
}

// GetByOrderID implements usecase.OrderRepo.
func (r *OrderRepo) GetByOrderID(ctx context.Context, number string) (entity.Order, error) {
	const sql = `SELECT id, user_id, status, accrual, uploaded_at
	FROM public.orders
	WHERE id = $1;`

	var err error
	var order entity.Order
	err = r.Pool.QueryRow(ctx, sql, number).Scan(&order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

// GetByStatus implements usecase.OrderRepo.
func (r *OrderRepo) GetByStatus(ctx context.Context, status []entity.Status) ([]entity.Order, error) {
	const sql = `SELECT id, user_id, status, accrual, uploaded_at
	FROM orders
	WHERE status = ANY($1);`

	var orders []entity.Order
	rows, err := r.Pool.Query(ctx, sql, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// GetByUserID implements usecase.OrderRepo.
func (r *OrderRepo) GetByUserID(ctx context.Context, userID int) ([]entity.Order, error) {
	const sql = `SELECT id, user_id, status, accrual, uploaded_at
	FROM orders
	WHERE user_id = $1;`

	var orders []entity.Order
	rows, err := r.Pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Update implements usecase.OrderRepo.
func (r *OrderRepo) Update(ctx context.Context, order entity.Order) error {
	const (
		sqlOrder   = "UPDATE orders SET status=$1, accrual=$2 WHERE id = $3;"
		sqlUserID  = "SELECT user_id FROM orders WHERE id = $1;"
		sqlBalance = "UPDATE balances SET current=current+$1 WHERE user_id=$2;"
	)

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sqlOrder, order.Status, order.Accrual, order.Number)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("can not update order error: %v", err)
	}

	var userID int
	err = tx.QueryRow(ctx, sqlUserID, order.Number).Scan(&userID)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("can not query userID error: %v", err)
	}

	_, err = tx.Exec(ctx, sqlBalance, order.Accrual, userID)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("can not update balance error: %v", err)
	}

	return tx.Commit(ctx)
}
