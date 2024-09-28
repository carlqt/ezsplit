package graph

import (
	"github.com/carlqt/ezsplit/config"
	"github.com/carlqt/ezsplit/internal/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Config       *config.EnvConfig
	Repositories *repository.Repository
}
