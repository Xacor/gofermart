package webapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Xacor/gophermart/internal/entity"
	"github.com/Xacor/gophermart/internal/utils/converter"
)

const (
	getOrderPath = "/api/orders/"
)

type AccrualsAPI struct {
	address string
	client  *http.Client
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

func NewAccrualsAPI(addr string, c *http.Client) *AccrualsAPI {
	return &AccrualsAPI{
		address: addr,
		client:  c,
	}
}

func (a *AccrualsAPI) GetOrderAccrual(ctx context.Context, number string) (*AccrualResponse, error) {
	url, err := a.orderPath(number)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		timeout, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			return nil, err
		}
		return nil, NewToManyRequestsError(resp.StatusCode, timeout)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var accrual AccrualResponse
	err = json.Unmarshal(body, &accrual)
	if err != nil {
		return nil, err
	}

	return &accrual, nil
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

func (a *AccrualsAPI) orderPath(number string) (string, error) {
	return url.JoinPath(a.address, getOrderPath, number)
}
