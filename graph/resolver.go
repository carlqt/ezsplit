package graph

import (
	"github.com/carlqt/ezsplit/internal"
	"github.com/carlqt/ezsplit/internal/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Config       *internal.EnvConfig
	Repositories *repository.Repository
}
