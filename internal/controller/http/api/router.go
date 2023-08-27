package api

import (
	"github.com/Xacor/gophermart/internal/controller/usecase"
	"github.com/Xacor/gophermart/internal/utils/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewRouter(handler chi.Router, l *zap.Logger, auth usecase.Auth, orders usecase.Orderer, balance usecase.Balancer, withdrawals usecase.Withdrawer, signKey string) {
	handler.Use(middleware.Logger)
	handler.Use(middleware.Compress(5, "application/json"))
	handler.Use(middleware.Recoverer)

	handler.Route("/api/user", func(r chi.Router) {
		newAuthRoutes(r, auth, l)
		handler.Route("/orders", func(r chi.Router) {
			r.Use(jwt.WithJWTAuth(signKey))
			newOrdersRoutes(r, orders, l)
		})
		handler.Route("/balance", func(r chi.Router) {
			r.Use(jwt.WithJWTAuth(signKey))
			newBalanceRoutes(r, balance, l)
		})

		handler.Route("/", func(r chi.Router) {
			r.Use(jwt.WithJWTAuth(signKey))
			newWithdrawalsRoutes(r, withdrawals, l)
		})

	})
}
