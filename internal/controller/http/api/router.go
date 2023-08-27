package api

import (
	"github.com/Xacor/gophermart/internal/controller/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewRouter(handler chi.Router, l *zap.Logger, auth usecase.Auth, orders usecase.Orderer, signKey string) {
	handler.Use(middleware.Logger) // TEMPRORARY SOLUTION
	handler.Use(middleware.Compress(5, "application/json"))
	handler.Use(middleware.Recoverer)

	handler.Route("/api/user", func(r chi.Router) {
		newAuthRoutes(r, auth, l)
		newOrdersRoutes(r, orders, l, signKey)

		// r.Route("/balance", func(r chi.Router) {
		// 	r.Get("/", nil)
		// 	r.Post("/withdraw", nil)
		// })

		// r.Get("/withdrawals", nil)
	})
}
