package product

import (
	"fmt"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	tests := []struct{
		name string
		ctx context.Context
		want []string
		wantErr bool
	}{{
		name: "get filenames",
		ctx: context.Background(),
		want: []string{"aaa", "bbb", "ccc"},
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
	type args struct{
		ctx context.Context
		filename string
	}

	tests := []struct{
		name string
		args args
		want error
	}{{
		name: "delete by filename",
		args: args{
			ctx: context.Background(),
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
