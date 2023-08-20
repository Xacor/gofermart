package entity

import "time"

type Status string

const (
	New        Status = "NEW"
	Invalid    Status = "INVALID"
	Processing Status = "PROCESSING"
	Processed  Status = "PROCESSED"
)

type Order struct {
	Number     string    `json:"number,omitempty"`
	UserID     int       `json:"-"`
	Status     Status    `json:"status,omitempty"`
	Accrual    int       `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at,omitempty"`
}
