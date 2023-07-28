//go:build integration

// *build integration

package integration

import (
	"crawler/murauchi"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestGetAllCategories(t *testing.T) {
	type want struct{
		err error
		categoriesCount int
	}
	tests := []struct{
		name string
		want want
	}{{
		name: "get all categories",
		want: want{nil, 1},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categories, err := murauchi.GetAllCategories()

			assert.GreaterOrEqual(t, len(categories), tt.want.categoriesCount)
			assert.NoError(t, err)
		})
	}
}
