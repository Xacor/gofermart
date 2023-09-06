package usecase

import (
	"context"

	"github.com/Xacor/gophermart/internal/entity"
	"github.com/jackc/pgx/v5"
)

// TODO: везде в repo передавать транзакции как аргумент
type (
	Txer interface {
		Begin(ctx context.Context) (pgx.Tx, error)
		Commit(ctx context.Context, tx pgx.Tx) error
		Rollback(ctx context.Context, tx pgx.Tx) error
	}

	Auth interface {
		Register(ctx context.Context, user entity.User) (string, error)
		Authenticate(ctx context.Context, user entity.User) (string, error)
	}

	UserRepo interface {
		Txer
		CreateUser(ctx context.Context, user entity.User) error
		GetByID(ctx context.Context, id int) (entity.User, error)
		GetByLogin(ctx context.Context, login string) (entity.User, error)
	}

	Orderer interface {
		GetOrders(ctx context.Context, userID int) ([]entity.Order, error)
		CreateOrder(ctx context.Context, number string, userID int) error
		PollOrders(ctx context.Context) error
	}

	OrderRepo interface {
		Txer
		Create(ctx context.Context, order entity.Order) error
		Update(ctx context.Context, order entity.Order) error // также должен обновлять кол-во бонусов
		GetByOrderID(ctx context.Context, number string) (entity.Order, error)
		GetByStatus(ctx context.Context, status []entity.Status) ([]entity.Order, error)
		GetByUserID(ctx context.Context, userID int) ([]entity.Order, error)
	}

	Balancer interface {
		GetUserBalance(ctx context.Context, userID int) (entity.Balance, error)
	}

	BalanceRepo interface {
		Txer
		Get(ctx context.Context, userID int) (entity.Balance, error)
	}

	Withdrawer interface {
		Withdraw(ctx context.Context, withdraw entity.Withdraw) error // также должен обновлять кол-во бонусов
		ListWithdrawals(ctx context.Context, userID int) ([]entity.Withdraw, error)
	}

	WithdrawalsRepo interface {
		Txer
		GetByUserID(ctx context.Context, userID int) ([]entity.Withdraw, error)
		Create(ctx context.Context, withdraw entity.Withdraw) error
	}
)
