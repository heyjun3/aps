package product

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"api-server/test"
)

func keepaSeed(repo KeepaRepository) error {
	keepas := make([]*Keepa, 100)
	for i := 0; i < 100; i++ {
		keepas[i] = &Keepa{Asin: "asin_" + strconv.Itoa(i)}
	}
	return repo.Save(context.Background(), keepas)
}

func TestKeepaGetByAsins(t *testing.T) {
	db := test.CreateTestDBConnection()
	if err := db.ResetModel(context.Background(), &Keepa{}); err != nil {
		panic(err)
	}
	repo := KeepaRepository{DB: db}
	if err := keepaSeed(repo); err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		args    []string
		want    []string
		wantErr bool
	}{
		{name: "keepas get by asins", args: []string{"asin_1"}, want: []string{"asin_1"}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ks, err := repo.GetByAsins(context.Background(), tt.args)

			assert.Equal(t, tt.want, ks.Asins())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestKeepaGetCounts(t *testing.T) {
	db := test.CreateTestDBConnection()
	if err := db.ResetModel(context.Background(), &Keepa{}); err != nil {
		panic(err)
	}
	ks := []*Keepa{
		{Asin: "aaa", Modified: time.Date(2022, 4, 1, 0, 0, 0, 0, time.Local)},
		{Asin: "bbb"},
		{Asin: "ccc"},
	}
	repo := KeepaRepository{DB: db}
	if err := repo.Save(context.Background(), ks); err != nil {
		panic(err)
	}

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
