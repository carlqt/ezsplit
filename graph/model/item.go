package model

type Item struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Price     *float64 `json:"price,omitempty"`
	ReceiptID string   `json:"receiptId"`
	SharedBy  []*User  `json:"sharedBy"`
}
