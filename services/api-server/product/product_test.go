package product

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func po[T any](v T) *T {
	return &v
}

func createTestData(db *bun.DB) {
	count := 150
	products := make([]Product, count)
	for i := 0; i < count; i++ {
		p := Product{Asin: "aaa" + fmt.Sprint(i), Filename: "aaa", Profit: po[int64](200),
			ProfitRate: po[float64](0.1), Unit: po[int64](1)}
		products[i] = p
	}
	keepas := make([]Keepa, count)
	for i := 0; i < count; i++ {
		k := Keepa{Asin: "aaa" + fmt.Sprint(i), Drops: 4 + i}
		keepas[i] = k
	}
	productRepo := ProductRepository{DB: db}
	keepaRepo := KeepaRepository{DB: db}
	if err := productRepo.Save(context.Background(), products); err != nil {
		panic("faild create test data")
	}
	if err := keepaRepo.Save(context.Background(), keepas); err != nil {
		panic("faild create test data")
	}
}

func TestGetCounts(t *testing.T) {
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		panic(fmt.Errorf("test database dsn is null"))
	}
	db := OpenDB(dsn)
	err := db.ResetModel(context.Background(), &Product{})
	if err != nil {
		panic(err)
	}
	p := []Product{
		{Asin: "aaa", Filename: "aaa", Price: po[int64](300)},
		{Asin: "bbb", Filename: "bbb", Price: po[int64](400), FeeRate: po[float64](0.1)},
		{Asin: "ccc", Filename: "ccc"},
	}
	repo := ProductRepository{DB: db}
	if err := repo.Save(context.Background(), p); err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		ctx     context.Context
		want    map[string]int
		wantErr bool
	}{
		{"get count", context.Background(), map[string]int{"total": 3, "price": 2, "fee": 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := repo.GetCounts(tt.ctx)

			assert.Equal(t, tt.want, count)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetFilenames(t *testing.T) {
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		panic(fmt.Errorf("test database dsn is null"))
	}
	db := OpenDB(dsn)
	err := db.ResetModel(context.Background(), &Product{})
	if err != nil {
		panic(err)
	}
	p := []Product{
		{Asin: "aaa", Filename: "aaa"},
		{Asin: "bbb", Filename: "bbb"},
		{Asin: "ccc", Filename: "ccc"},
	}
	repo := ProductRepository{DB: db}
	if err := repo.Save(context.Background(), p); err != nil {
		panic(err)
	}
	tests := []struct {
		name    string
		ctx     context.Context
		want    []string
		wantErr bool
	}{{
		name:    "get filenames",
		ctx:     context.Background(),
		want:    []string{"aaa", "bbb", "ccc"},
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filenames, err := repo.GetFilenames(tt.ctx)

			assert.Equal(t, tt.want, filenames)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetProductWithChart(t *testing.T) {
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		panic(fmt.Errorf("test database dsn is null"))
	}
	db := OpenDB(dsn)
	if err := db.ResetModel(context.Background(), (*Product)(nil)); err != nil {
		panic("test database dsn is null")
	}
	if err := db.ResetModel(context.Background(), (*Keepa)(nil)); err != nil {
		panic("test database dsn is null")
	}
	createTestData(db)

	type args struct {
		ctx      context.Context
		filename string
		page     int
		limit    int
	}
	type want struct {
		count int
		total int
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{"get product with chart", args{context.Background(), "aaa", 1, 100}, want{100, 150}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := ProductRepository{DB: db}
			ps, total, err := repo.GetProductWithChart(tt.args.ctx, tt.args.filename, tt.args.page, tt.args.limit)
			assert.Equal(t, tt.want.count, len(ps))
			assert.Equal(t, tt.want.total, total)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteByFilename(t *testing.T) {
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		panic(fmt.Errorf("test database dsn is null"))
	}
	db := OpenDB(dsn)
	err := db.ResetModel(context.Background(), &Product{})
	if err != nil {
		panic(err)
	}
	p := []Product{
		{Asin: "aaa", Filename: "aaa"},
		{Asin: "bbb", Filename: "bbb"},
		{Asin: "ccc", Filename: "ccc"},
	}
	repo := ProductRepository{DB: db}
	if err := repo.Save(context.Background(), p); err != nil {
		panic(err)
	}
	type args struct {
		ctx      context.Context
		filename string
	}

	tests := []struct {
		name string
		args args
		want error
	}{{
		name: "delete by filename",
		args: args{
			ctx:      context.Background(),
			filename: "aaa",
		},
		want: nil,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteByFilename(tt.args.ctx, tt.args.filename)

			assert.Equal(t, tt.want, err)
		})
	}
}
