package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"api-server/test"
)

func TestSavePrices(t *testing.T) {
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
			test.OmitErr(NewCurrentPrice(Ptr("sku"), Ptr(1), Ptr(2))),
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
