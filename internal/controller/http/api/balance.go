package api

import (
	"encoding/json"
	"net/http"

	"github.com/Xacor/gophermart/internal/controller/usecase"
	"github.com/Xacor/gophermart/internal/utils/converter"
	"github.com/Xacor/gophermart/internal/utils/jwt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type balanceRoutes struct {
	b usecase.Balancer
	l *zap.Logger
}

func newBalanceRoutes(r chi.Router, balancer usecase.Balancer, l *zap.Logger) {
	br := &balanceRoutes{balancer, l}
	r.Get("/", br.GetBalance)
}

type balanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (br *balanceRoutes) GetBalance(w http.ResponseWriter, r *http.Request) {

	userID := jwt.GetUserIDFromCtx(r.Context())
	balance, err := br.b.GetUserBalance(r.Context(), userID)
	if err != nil {
		br.l.Error("can not query user balance", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balanceResp := balanceResponse{
		Current:   converter.IntToFloat(balance.Current),
		Withdrawn: converter.IntToFloat(balance.Withdrawn),
	}

	br.l.Debug("GetBalance", zap.Any("resp", balanceResp))

	body, err := json.Marshal(balanceResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
