package graph

import (
	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ReceiptStore map[string]model.Receipt
	UserStore    map[string]model.User
	Repositories *repository.Repository
}
