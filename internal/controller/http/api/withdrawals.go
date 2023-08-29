package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/Xacor/gophermart/internal/controller/usecase"
	"github.com/Xacor/gophermart/internal/entity"
	"github.com/Xacor/gophermart/internal/utils/converter"
	"github.com/Xacor/gophermart/internal/utils/jwt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type withdrawalsRoutes struct {
	w usecase.Withdrawer
	l *zap.Logger
}

func newWithdrawalsRoutes(handler chi.Router, withdrawals usecase.Withdrawer, l *zap.Logger) {
	wr := &withdrawalsRoutes{withdrawals, l}
	l.Debug("newWithdrawalRoute")
	handler.Post("/balance/withdraw", http.HandlerFunc(wr.PostWithdraw))
	handler.Get("/withdrawals", wr.ListWithdrawals)
}

type withdraw struct {
	Order       string    `json:"order,omitempty"`
	Sum         float64   `json:"sum,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}

func (wr *withdrawalsRoutes) ListWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(jwt.UserIDKey).(int)
	withdrawals, err := wr.w.ListWithdrawals(r.Context(), userID)
	if err != nil {
		wr.l.Error("list withdrawals handler", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make([]withdraw, 0)
	for _, v := range withdrawals {
		resp = append(resp, withdraw{
			Order:       v.Order,
			Sum:         converter.IntToFloat(v.Sum),
			ProcessedAt: v.ProcessedAt,
		})
	}

	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
	w.Header().Set("Content-Type", "application/json")
}

func (wr *withdrawalsRoutes) PostWithdraw(w http.ResponseWriter, r *http.Request) {
	var request withdraw
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value(jwt.UserIDKey).(int)
	withdraw := entity.Withdraw{
		UserID: userID,
		Order:  request.Order,
		Sum:    converter.FloatToInt(request.Sum),
	}
	wr.l.Debug("widraw req", zap.Any("withdraw", withdraw), zap.Int("userID", withdraw.UserID))

	if goluhn.Validate(withdraw.Order) != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = wr.w.Withdraw(r.Context(), withdraw)
	if err != nil {
		if errors.Is(err, usecase.ErrInsufficientBalance) {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
