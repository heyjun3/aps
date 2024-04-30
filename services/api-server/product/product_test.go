package product

import (
	"api-server/test"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func Ptr[T any](v T) *T {
	return &v
}

func createTestData(db *bun.DB) {
	count := 150
	products := make([]Product, count)
	for i := 0; i < count; i++ {
		p := Product{Asin: "aaa" + fmt.Sprint(i), Filename: "aaa", Profit: Ptr[int64](200),
			ProfitRate: Ptr[float64](0.1), Unit: Ptr[int64](1)}
		products[i] = p
	}
	keepas := make([]*Keepa, count)
	for i := 0; i < count; i++ {
		k := Keepa{Asin: "aaa" + fmt.Sprint(i), Drops: 4 + i}
		keepas[i] = &k
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
	db := test.CreateTestDBConnection()
	err := db.ResetModel(context.Background(), &Product{})
	if err != nil {
		panic(err)
	}
	p := []Product{
		{Asin: "aaa", Filename: "aaa", Price: Ptr[int64](300)},
		{Asin: "bbb", Filename: "bbb", Price: Ptr[int64](400), FeeRate: Ptr[float64](0.1)},
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
	db := test.CreateTestDBConnection()
	err := db.ResetModel(context.Background(), &Product{})
	if err != nil {
		panic(err)
	}
	p := []Product{
		{Asin: "aaa", Filename: "aaa", Profit: Ptr[int64](200)},
		{Asin: "bbb", Filename: "bbb", Profit: Ptr[int64](200)},
		{Asin: "ccc", Filename: "ccc", Profit: Ptr[int64](200)},
		{Asin: "ddd1", Filename: "ddd", Profit: Ptr[int64](200)},
		{Asin: "ddd2", Filename: "ddd"},
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
	setupData := func() *bun.DB {
		db := test.CreateTestDBConnection()
		if err := db.ResetModel(context.Background(), (*Product)(nil)); err != nil {
			panic("test database dsn is null")
		}
		if err := db.ResetModel(context.Background(), (*Keepa)(nil)); err != nil {
			panic("test database dsn is null")
		}
		createTestData(db)
		return db
	}

	testGetProductWithChartBySearchCondition := func(t *testing.T, db *bun.DB) {
		repo := ProductRepository{DB: db}
		c := NewSearchCondition("aaa")
		ps, total, err := repo.GetProductWithChartBySearchCondition(
			context.Background(), c)
		assert.Equal(t, 100, len(ps))
		assert.Equal(t, 150, total)
		assert.NoError(t, err)
	}

	tests := []struct {
		name string
		fn   func(*testing.T, *bun.DB)
	}{
		{"get product by search condition", testGetProductWithChartBySearchCondition},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupData()
			tt.fn(t, db)
		})
	}
}

func TestDeleteByFilename(t *testing.T) {
	db := test.CreateTestDBConnection()
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

func TestDeleteIfCondition(t *testing.T) {
	db := test.CreateTestDBConnection()
	err := db.ResetModel(context.Background(), &Product{})
	if err != nil {
		panic(err)
	}
	now := time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local)
	products := []Product{
		{Asin: "test", Filename: "test", CreatedAt: now},
		{Asin: "test1", Filename: "test", Profit: Ptr(int64(199))},
		{Asin: "test2", Filename: "test", ProfitRate: Ptr(float64(0.09))},
		{Asin: "test3", Filename: "test", Unit: Ptr(int64(3))},
	}
	repo := ProductRepository{DB: db}
	if err := repo.Save(context.Background(), products); err != nil {
		panic(err)
	}

	type args struct {
		ctx       context.Context
		condition Condition
	}
	type want struct {
		count    int
		products []Product
	}
	tests := []struct {
		name  string
		args  args
		want  want
		isErr bool
	}{{
		name: "delete rows if condition",
		args: args{
			ctx:       context.Background(),
			condition: *NewCondition(200, 2, 0.1),
		},
		want: want{
			count:    1,
			products: []Product{{Asin: "test", Filename: "test", CreatedAt: now}},
		},
		isErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteIfCondition(tt.args.ctx, &tt.args.condition)

			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			products, err := repo.GetProduct(context.Background())
			if err != nil {
				panic(err)
			}
			assert.Equal(t, tt.want.count, len(products))
			assert.Equal(t, tt.want.products, products)
		})
	}
}

func TestDeleteIfConditionWithKeepa(t *testing.T) {
	db := test.CreateTestDBConnection()
	for _, model := range []interface{}{(*Product)(nil), (*Keepa)(nil)} {
		if err := db.ResetModel(context.Background(), model); err != nil {
			panic(err)
		}
	}
	now := time.Date(2022, 2, 2, 0, 0, 0, 0, time.Local)
	products := []Product{
		{Asin: "test", Filename: "test", CreatedAt: now},
		{Asin: "test1", Filename: "test", CreatedAt: now},
		{Asin: "test2", Filename: "test", CreatedAt: now},
		{Asin: "test3", Filename: "test", CreatedAt: now},
	}
	repo := ProductRepository{DB: db}
	if err := repo.Save(context.Background(), products); err != nil {
		panic(err)
	}

	keepas := []*Keepa{
		{Asin: "test", Drops: 3},
		{Asin: "test1", Drops: 4},
		{Asin: "test2", Drops: 5},
	}
	if err := (KeepaRepository{DB: db}).Save(context.Background(), keepas); err != nil {
		panic(err)
	}

	type want struct {
		count    int
		products []Product
	}
	tests := []struct {
		name  string
		want  want
		isErr bool
	}{{
		name: "delete row if condition",
		want: want{
			count: 3,
			products: []Product{
				{Asin: "test1", Filename: "test", CreatedAt: now},
				{Asin: "test2", Filename: "test", CreatedAt: now},
				{Asin: "test3", Filename: "test", CreatedAt: now},
			},
		},
		isErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteIfConditionWithKeepa(context.Background())
			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			products, err := repo.GetProduct(context.Background())
			if err != nil {
				panic(err)
			}
			assert.Equal(t, tt.want.count, len(products))
			assert.Equal(t, tt.want.products, products)
		})
	}
}

func TestRefreshGeneratedColumns(t *testing.T) {
	db := test.CreateTestDBConnection()
	t.Run("test", func(t *testing.T) {
		err := ProductRepository{DB: db}.RefreshGeneratedColumns(context.Background())

		assert.NoError(t, err)
	})
}
