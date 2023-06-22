package product

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKeepaGetCounts(t *testing.T) {
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		panic(fmt.Errorf("test database dsn is null"))
	}
	db := OpenDB(dsn)
	err := db.ResetModel(context.Background(), &Keepa{})
	if err != nil {
		panic(err)
	}
	ks := []Keepa{
		{Asin: "aaa", Modified: time.Date(2022, 4, 1, 0, 0, 0, 0, time.Local)},
		{Asin: "bbb"},
		{Asin: "ccc"},
	}
	repo := KeepaRepository{DB: db}
	err = repo.Save(context.Background(), ks)
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Now())
	tests := []struct {
		name    string
		ctx     context.Context
		want    map[string]int
		wantErr bool
	}{
		{"get counts", context.Background(), map[string]int{"total": 3, "modified": 2}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetCounts(tt.ctx)

			assert.Equal(t, tt.want, result)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
