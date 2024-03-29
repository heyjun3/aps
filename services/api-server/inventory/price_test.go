package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"api-server/test"
)

func TestSaveCurrentPrices(t *testing.T) {
	ctx := context.Background()
	db := test.CreateTestDBConnection()
	test.ResetModel(ctx, db, &CurrentPrice{})
	repo := PriceRepository[*CurrentPrice]{}

	tests := []struct {
		name    string
		prices  CurrentPrices
		wantErr bool
	}{{
		name: "save current prices",
		prices: CurrentPrices{
			test.OmitErr(NewCurrentPrice("sku", Ptr(1.0), Ptr(2.0))),
			test.OmitErr(NewCurrentPrice("skux", Ptr(1.0), Ptr(2.0))),
		},
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Save(ctx, db, tt.prices)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSaveLowestPrices(t *testing.T) {
	ctx := context.Background()
	db := test.CreateTestDBConnection()
	test.ResetModel(ctx, db, &CurrentPrice{})
	repo := PriceRepository[*LowestPrice]{}

	tests := []struct {
		name    string
		prices  LowestPrices
		wantErr bool
	}{{
		name: "save current prices",
		prices: LowestPrices{
			test.OmitErr(NewLowestPrice("sku", Ptr(1.0), Ptr(2.0))),
			test.OmitErr(NewLowestPrice("sku1", Ptr(1.0), Ptr(2.0))),
		},
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Save(ctx, db, tt.prices)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSaveDesiredPrice(t *testing.T) {
	db := test.CreateTestDBConnection()
	test.ResetModel(context.Background(), db, &DesiredPrice{})
	repo := PriceRepository[*DesiredPrice]{}

	tests := []struct {
		name    string
		prices  DesiredPrices
		wantErr bool
	}{{
		name: "save desired prices",
		prices: DesiredPrices{
			test.OmitErr(NewDesiredPrice(Ptr("sku1"), Ptr(10000.0), Ptr(1.0), *test.OmitErr(NewLowestPrice("sku1", Ptr(10000.0), Ptr(2.0))))),
			test.OmitErr(NewDesiredPrice(Ptr("sku2"), Ptr(20000.0), Ptr(2.0), *test.OmitErr(NewLowestPrice("sku2", Ptr(20000.0), Ptr(4.0))))),
		},
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Save(context.Background(), db, tt.prices)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
