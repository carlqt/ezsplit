package model

type Item struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Price     string  `json:"price"`
	ReceiptID string  `json:"receiptId"`
	SharedBy  []*User `json:"sharedBy"`
}
