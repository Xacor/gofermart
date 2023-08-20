package entity

type Balance struct {
	UserID    int `json:"-"`
	Current   int `json:"current,omitempty"`
	Withdrawn int `json:"withdrawn,omitempty"`
}
