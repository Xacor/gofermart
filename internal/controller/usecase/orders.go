package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/Xacor/gophermart/internal/controller/usecase/webapi"
	"github.com/Xacor/gophermart/internal/entity"
	"github.com/Xacor/gophermart/internal/utils/converter"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const pollInterval = time.Second * 3

var (
	ErrAnothersOrder  = errors.New("another user's order")
	ErrInvalidLuhn    = errors.New("invalid luhn number")
	ErrAlredyUploaded = errors.New("order already uploaded by this user")
)

type OrderUseCase struct {
	orderRepo OrderRepo
	api       *webapi.AccrualsAPI
	l         *zap.Logger
}

func NewOrdersUseCase(ctx context.Context, orderRepo OrderRepo, api *webapi.AccrualsAPI, logger *zap.Logger) *OrderUseCase {
	usecase := &OrderUseCase{
		orderRepo: orderRepo,
		api:       api,
		l:         logger,
	}
	go usecase.PollOrders(ctx)

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

	tx, err := o.orderRepo.Begin(ctx)
	if err != nil {
		return err
	}

	uploaded, err := o.orderRepo.GetByOrderID(ctx, number, tx)
	if err == nil {
		tx.Rollback(ctx)
		if uploaded.UserID == userID {
			return ErrAlredyUploaded // 200
		} else {
			return ErrAnothersOrder //409
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	order := entity.Order{
		Number:     number,
		UserID:     userID,
		Status:     entity.New,
		UploadedAt: time.Now(),
	}

	err = o.orderRepo.Create(ctx, order, tx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

// Полит внешнее api на наличие бонусов по заказу и сохраняет начисления в бд
func (o *OrderUseCase) PollOrders(ctx context.Context) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := o.queryAndUpdate(ctx)

			if err != nil {
				o.l.Error("error update orders status", zap.Error(err))
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

	g, groupCtx := errgroup.WithContext(ctx)

	for _, order := range orders {
		order := order

		g.Go(func() error {
			resp, err := o.api.GetOrderAccrual(groupCtx, order.Number)
			if err != nil {
				return fmt.Errorf("api error: %v", err)
			}

			accrualToOrder(resp, &order)
			err = o.orderRepo.Update(groupCtx, order)
			if err != nil {
				return err
			}

			return nil
		})
	}

	return g.Wait()
}

func accrualToOrder(accrual *webapi.AccrualResponse, order *entity.Order) {
	const (
		statusRegistered = "REGISTERED"
		statusInvalid    = "INVALID"
		statusProcessing = "PROCESSING"
		statusProcessed  = "PROCESSED"
	)

	if accrual.Order != order.Number {
		return
	}

	switch accrual.Status {
	case statusRegistered:
		order.Status = entity.New
	case statusInvalid:
		order.Status = entity.Invalid
	case statusProcessing:
		order.Status = entity.Processing
	case statusProcessed:
		order.Status = entity.Processed
		order.Accrual = converter.FloatToInt(accrual.Accrual)
	}
}
