package scrape

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"crawler/test/util"
)

func TestRunServiceStatusSave(t *testing.T) {
	db, ctx := util.DatabaseFactory()
	db.ResetModel(ctx, (*RunServiceHistory)(nil))
	repo := RunServiceHistoryRepository{}
	history := NewRunServiceHistory("yahoo", "https://yahoo.co.jp", "DONE")

	newHistory, err := repo.Save(ctx, db, history)

	assert.NoError(t, err)
	assert.Equal(t, history, newHistory)
}
