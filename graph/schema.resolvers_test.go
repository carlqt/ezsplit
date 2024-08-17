package graph

import (
	"context"
	"testing"

	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal"
	integration_test "github.com/carlqt/ezsplit/test"
	"github.com/stretchr/testify/assert"
)

func TestSchemaResolver(t *testing.T) {
	app := internal.NewApp()
	resolvers := Resolver{Repositories: app.Repositories, Config: app.Config}
	testItemResolver := itemResolver{&resolvers}
	truncateTables := func() {
		integration_test.TruncateAllTables(app.DB)
	}

	t.Run("SharedBy", func(t *testing.T) {
		t.Run("when there are no results", func(t *testing.T) {
			defer truncateTables()

			modelItem := model.Item{ID: "888"}
			result, _ := testItemResolver.SharedBy(context.TODO(), &modelItem)

			assert.Empty(t, result, result)
			assert.Zero(t, len(result))
		})
	})
}
