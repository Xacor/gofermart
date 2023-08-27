package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Xacor/gophermart/internal/controller/usecase"
	"github.com/Xacor/gophermart/internal/utils/jwt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ordersRoutes struct {
	o usecase.Orderer
	l *zap.Logger
}

func newOrdersRoutes(handler chi.Router, o usecase.Orderer, l *zap.Logger, signKey string) {
	or := &ordersRoutes{o, l}
	handler.Route("/orders", func(r chi.Router) {
		r.Use(jwt.WithJWTAuth(signKey))
		r.Post("/", or.PostOrder)
		r.Get("/", or.GetOrders)
	})

}

func (or *ordersRoutes) PostOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		or.l.Error("error read body", zap.Error(err), zap.Int("body length", len(body)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	number := string(body)
	userID := r.Context().Value(jwt.UserIDKey).(int)

	err = or.o.CreateOrder(r.Context(), number, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidLuhn) {
			or.l.Error("error create order", zap.Error(err), zap.String("number", number))
			w.WriteHeader(http.StatusUnprocessableEntity)
			return

		} else if errors.Is(err, usecase.ErrAnothersOrder) {
			or.l.Error("error create order", zap.Error(err))
			w.WriteHeader(http.StatusConflict)
			return

		} else if errors.Is(err, usecase.ErrAlredyUploaded) {
			or.l.Error("error create order", zap.Error(err))
			w.WriteHeader(http.StatusOK)
			return

		}
		or.l.Error("another error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (or *ordersRoutes) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(jwt.UserIDKey).(int)
	orders, err := or.o.GetOrders(r.Context(), userID)
	if err != nil {
		or.l.Error("error get orders", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	json, err := json.Marshal(orders)
	if err != nil {
		or.l.Error("error marshling orders", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(json)
}
