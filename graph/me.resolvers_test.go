package graph

import (
	"context"
	"testing"

	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal"
	"github.com/stretchr/testify/assert"
)

func newMeResolver() meResolver {
	app := internal.NewApp()
	resolvers := Resolver{Repositories: app.Repositories, Config: app.Config}

	return meResolver{&resolvers}
}

func TestTotalPayables(t *testing.T) {
	r := newMeResolver()
	ctx := context.TODO()
	meModel := &model.Me{ID: "1"}

	t.Run("when there are no rows in user_orders", func(t *testing.T) {
		result, err := r.TotalPayables(ctx, meModel)

		if assert.Nil(t, err) {
			assert.Equal(t, "0.00", result)
		}
	})
}
