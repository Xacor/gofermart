package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/Xacor/gophermart/internal/controller/usecase/webapi"
	"github.com/Xacor/gophermart/internal/entity"
	"go.uber.org/zap"
)

const pollInterval = time.Second * 3

var (
	ErrAnothersOrder  = errors.New("another user's order")
	ErrInvalidLuhn    = errors.New("invalid luhn number")
	ErrAlredyUploaded = errors.New("order already uploaded by this user")
)

type OrderUseCase struct {
	orderRepo   OrderRepo
	balanceRepo BalanceRepo
	api         *webapi.AccrualsAPI
	l           *zap.Logger
}

func NewOrdersUseCase(orderRepo OrderRepo, balanceRepo BalanceRepo, api *webapi.AccrualsAPI, logger *zap.Logger) *OrderUseCase {
	usecase := &OrderUseCase{
		orderRepo:   orderRepo,
		balanceRepo: balanceRepo,
		api:         api,
		l:           logger,
	}
	go usecase.PollOrders(context.Background())

	return usecase
}

func (o *OrderUseCase) GetOrders(ctx context.Context, userID int) ([]entity.Order, error) {
	return o.orderRepo.GetByUserID(ctx, userID)
}

func (o *OrderUseCase) CreateOrder(ctx context.Context, number string, userID int) error {
	err := goluhn.Validate(number)
	if err != nil {
		return ErrInvalidLuhn // 422
	}

	uploaded, err := o.orderRepo.GetByOrderID(ctx, number)
	if err == nil {
		if uploaded.UserID == userID {
			return ErrAlredyUploaded // 200
		} else {
			return ErrAnothersOrder //409
		}
	}

	order := entity.Order{
		Number:     number,
		UserID:     userID,
		Status:     entity.New,
		UploadedAt: time.Now(),
	}

	return o.orderRepo.Create(ctx, order)
}

// Полит внешнее api на наличие бонусов по заказу и сохраняет начисления в бд
func (o *OrderUseCase) PollOrders(ctx context.Context) error {
	o.l.Debug("polling started")
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			o.l.Debug("tick")
			err := o.queryAndUpdate(ctx)

			if err != nil {
				o.l.Error("error update orders status", zap.Error(err))
				var reqError *webapi.ToManyRequestsError
				if errors.As(err, &reqError) {
					<-time.After(time.Duration(reqError.RetryAfter) * time.Second)
				}
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (o *OrderUseCase) queryAndUpdate(ctx context.Context) error {
	orders, err := o.orderRepo.GetByStatus(context.Background(), []entity.Status{entity.New, entity.Processing})
	if err != nil {
		return err
	}

	if len(orders) == 0 {
		return nil
	}

	for _, order := range orders {
		resp, err := o.api.GetOrderAccrual(ctx, order.Number)
		if err != nil {
			return fmt.Errorf("api error", err)
		}

		o.api.AccrualToOrder(resp, &order)
		err = o.orderRepo.Update(ctx, order)
		if err != nil {
			return err
		}
	}

	return nil
}
