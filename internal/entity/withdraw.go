package entity

import "time"

type Withdraw struct {
	Order       string    `json:"order,omitempty"`
	Sum         int       `json:"sum,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}
