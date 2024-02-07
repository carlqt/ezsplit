package model

type Receipt struct {
	ID          string  `json:"id"`
	User        *User   `json:"ownedBy,omitempty"`
	Description string  `json:"description"`
	Total       string  `json:"total"`
	Items       []*Item `json:"items"`
	UserID      string  `json:"userId"`
}
