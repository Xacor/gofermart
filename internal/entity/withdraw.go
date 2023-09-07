package entity

import "time"

type Withdraw struct {
	ID          int       `json:"-"`
	UserID      int       `json:"-"`
	Order       string    `json:"order,omitempty"`
	Sum         int       `json:"sum,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}
