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
			test.OmitErr(NewCurrentPrice("sku", 1, 2)),
			test.OmitErr(NewCurrentPrice("skux", 1, 2)),
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
			test.OmitErr(NewLowestPrice("sku", 1, 2)),
			test.OmitErr(NewLowestPrice("sku1", 1, 2)),
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
