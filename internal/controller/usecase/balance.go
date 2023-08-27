package usecase

import (
	"context"

	"github.com/Xacor/gophermart/internal/entity"
	"go.uber.org/zap"
)

type BalanceUseCase struct {
	repo BalanceRepo
	l    *zap.Logger
}

func NewBalanceUseCase(repo BalanceRepo, logger *zap.Logger) *BalanceUseCase {
	return &BalanceUseCase{repo, logger}
}

func (b *BalanceUseCase) GetUserBalance(ctx context.Context, userID int) (entity.Balance, error) {
	return b.repo.Get(ctx, userID)
}
