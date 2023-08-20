package webapi

import (
	"errors"
)

var ErrToManyRequests = errors.New("to many requests")

type ToManyRequestsError struct {
	StatusCode int
	RetryAfter int
	Err        error
}

func (e *ToManyRequestsError) Error() string {
	return e.Err.Error()
}

func NewToManyRequestsError(statusCode int, timeout int) error {
	return &ToManyRequestsError{
		StatusCode: statusCode,
		RetryAfter: timeout,
		Err:        ErrToManyRequests,
	}
}
