package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Xacor/gophermart/internal/config"
	"github.com/Xacor/gophermart/internal/controller/http/api"
	"github.com/Xacor/gophermart/internal/controller/usecase"
	repo "github.com/Xacor/gophermart/internal/controller/usecase/repo/postgres"
	"github.com/Xacor/gophermart/pkg/httpserver"
	"github.com/Xacor/gophermart/pkg/logger"
	"github.com/Xacor/gophermart/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.LogLevel)

	migrate(cfg.DatabaseURI, l)
	pg, err := postgres.New(cfg.DatabaseURI)
	if err != nil {
		l.Fatal("failed to init DB", zap.Error(err))
	}
	defer pg.Close()

	handler := chi.NewMux()

	auth := usecase.NewAuthUseCase(repo.NewUserRepo(pg), cfg.SecretKey)

	api.NewRouter(handler, l, auth)

	l.Info("starting HTTP server", zap.String("addr", cfg.Address))
	httpServer := httpserver.New(handler, httpserver.Address(cfg.Address))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("shutting down gracefully", zap.String("signal", s.String()))
	case err := <-httpServer.Notify():
		l.Error("httpServer failed to start", zap.Error(err))
	}

	if err := httpServer.Shutdown(); err != nil {
		l.Error("failed to shutdown httpServer", zap.Error(err))
	}
}
