package model

type Receipt struct {
	ID          string   `json:"id"`
	User        *User    `json:"ownedBy,omitempty"`
	Description string   `json:"description"`
	Total       *float64 `json:"total,omitempty"`
	Items       []*Item  `json:"items"`
	UserID      string   `json:"userId"`
}
