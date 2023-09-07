package webapi

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

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

func NewAccrualsAPI(addr string, c *resty.Client) *AccrualsAPI {
	return &AccrualsAPI{
		client: c.SetBaseURL(addr).
			SetRetryAfter(retryOnToManyRequests).
			AddRetryCondition(toManyRequestCond),
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

func retryOnToManyRequests(c *resty.Client, r *resty.Response) (time.Duration, error) {
	// обработка по дефолтной схеме
	if r.StatusCode() != http.StatusTooManyRequests {
		return 0, nil
	}

	timeout, err := strconv.Atoi(r.Header().Get("Retry-After"))
	if err != nil {
		return 0, err
	}
	return time.Second * time.Duration(timeout), nil
}

func toManyRequestCond(r *resty.Response, err error) bool {
	return r.StatusCode() == http.StatusTooManyRequests
}
