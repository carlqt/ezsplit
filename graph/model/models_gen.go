// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AddItemToReceiptInput struct {
	ReceiptID string   `json:"receiptId"`
	Name      string   `json:"name"`
	Price     *float64 `json:"price,omitempty"`
}

type AssignOrDeleteMeToItemInput struct {
	ItemID string `json:"itemId"`
}

type AssignUserToItemInput struct {
	ItemID string `json:"itemId"`
	UserID string `json:"userId"`
}

type DeleteItemPayload struct {
	Msg string `json:"msg"`
	ID  string `json:"id"`
}

type DeleteMyReceiptInput struct {
	ID string `json:"id"`
}

type Item struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    string  `json:"price"`
	SharedBy []*User `json:"sharedBy"`
}

type LoginUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Me struct {
	ID            string     `json:"id"`
	Username      string     `json:"username"`
	TotalPayables string     `json:"totalPayables"`
	Receipts      []*Receipt `json:"receipts"`
}

type Mutation struct {
}

type Query struct {
}

type Receipt struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userId"`
	User        *User   `json:"user,omitempty"`
	Description string  `json:"description"`
	Total       string  `json:"total"`
	Slug        string  `json:"slug"`
	Items       []*Item `json:"items"`
}

type ReceiptInput struct {
	Description string   `json:"description"`
	Total       *float64 `json:"total,omitempty"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UserInput struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type UserWithJwt struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
}
