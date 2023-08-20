package postgres

import (
	"errors"

	"github.com/Xacor/gophermart/pkg/postgres"
)

var (
	ErrNoRows = errors.New("no rows")
)

type OrderRepo struct {
	*postgres.Postgres
}

func NewOrderRepo(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{pg}
}
