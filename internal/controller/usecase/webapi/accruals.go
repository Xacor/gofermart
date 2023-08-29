package webapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/Xacor/gophermart/internal/entity"
	"github.com/Xacor/gophermart/internal/utils/converter"
	"github.com/go-resty/resty/v2"
)

type AccrualsAPI struct {
	client *resty.Client
}

type AccrualResponse struct {
	Order   string  `json:"order,omitempty"`
	Status  string  `json:"status,omitempty"`
	Accrual float64 `json:"accrual,omitempty"`
}

const (
	statusRegistered = "REGISTERED"
	statusInvalid    = "INVALID"
	statusProcessing = "PROCESSING"
	statusProcessed  = "PROCESSED"
)

func NewAccrualsAPI(addr string, c *resty.Client) *AccrualsAPI {
	return &AccrualsAPI{
		client: c.SetBaseURL("http://" + addr),
	}
}

func (a *AccrualsAPI) GetOrderAccrual(ctx context.Context, number string) (*AccrualResponse, error) {
	var acc AccrualResponse
	req := a.client.R().
		SetResult(&acc).
		SetPathParam("number", number)

	resp, err := req.Get("/api/orders/{number}")
	if err != nil {
		return &AccrualResponse{}, err
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil, errors.New("no content")
	}

	if resp.StatusCode() == http.StatusTooManyRequests {
		return nil, errors.New("to many requests")
	}

	return &acc, nil

}

func (a *AccrualsAPI) AccrualToOrder(accrual *AccrualResponse, order *entity.Order) {
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
