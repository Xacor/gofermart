package usecase

import (
	"context"

	"github.com/Xacor/gophermart/internal/entity"
)

type (
	Auth interface {
		Register(ctx context.Context, user entity.User) error
		Authenticate(ctx context.Context, user entity.User) (string, error)
	}

	UserRepo interface {
		CreateUser(ctx context.Context, user entity.User) error
		GetByID(ctx context.Context, id string) (entity.User, error)
		GetByLogin(ctx context.Context, login string) (entity.User, error)
	}

	Orderer interface {
		GetOrders(ctx context.Context, userID int) ([]entity.Order, error)
		CreateOrder(ctx context.Context, number string, userID int) error
		PollOrders()
	}

	OrderRepo interface {
		Create(ctx context.Context, order entity.Order) error
		Get(ctx context.Context, order entity.Order) (entity.Order, error)
		Update(ctx context.Context, order entity.Order) error // также должен обновлять кол-во бонусов
		GetByOrderID(ctx context.Context, number string) (entity.Order, error)
		GetByStatus(ctx context.Context, status []entity.Status) ([]entity.Order, error)
		GetByUserID(ctx context.Context, userID int) ([]entity.Order, error) // возвращает пустой сдайс и ошиюку, если не найдены записи
	}

	BalanceRepo interface {
		Get(ctx context.Context, userID int) (entity.Balance, error)
		AddBonuses(ctx context.Context, userID, amount int) error
		Withdraw(ctx context.Context, userID, amount int) error
	}
)
