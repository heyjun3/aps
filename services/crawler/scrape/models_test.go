package scrape

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func NewTestProduct(name, productCode, url, jan, shopCode string, price int64) *Product {
// 	p, _ := NewProduct(name, productCode, url, jan, shopCode, price)
// 	return p
// }

// func TestMappingProducts(t *testing.T) {
// 	type args struct {
// 		mergeProducts  Products
// 		targetProducts Products
// 	}
// 	type want struct {
// 		mergedProducts Products
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want want
// 	}{
// 		{
// 			name: "merge products",
// 			args: args{
// 				mergeProducts: Products{
// 					(NewTestProduct("test", "test", "http://test.jp", "", "test_shop", 1111)),
// 					(NewTestProduct("test1", "test1", "http://test.jp", "", "test_shop", 1111)),
// 					(NewTestProduct("test2", "test2", "http://test.jp", "", "test_shop", 1111)),
// 				},
// 				targetProducts: Products{
// 					(NewTestProduct("test", "test", "test", "4444", "test_shop", 4444)),
// 					(NewTestProduct("test", "test1", "test1", "555", "test_shop", 4444)),
// 					(NewTestProduct("test", "test2", "test2", "7777", "test_shop", 4444)),
// 				},
// 			},
// 			want: want{
// 				mergedProducts: Products{
// 					(NewTestProduct("test", "test", "http://test.jp", "4444", "test_shop", 1111)),
// 					(NewTestProduct("test1", "test1", "http://test.jp", "555", "test_shop", 1111)),
// 					(NewTestProduct("test2", "test2", "http://test.jp", "7777", "test_shop", 1111)),
// 				},
// 			},
// 		},
// 		{
// 			name: "merge product is empty",
// 			args: args{
// 				mergeProducts: Products{},
// 				targetProducts: Products{
// 					(NewTestProduct("test", "test", "test", "11111", "test_shop", 4444)),
// 					(NewTestProduct("test", "test", "test1", "55555", "test_shop", 4444)),
// 				},
// 			},
// 			want: want{mergedProducts: Products{}},
// 		},
// 		{
// 			name: "target product is empty",
// 			args: args{
// 				mergeProducts: Products{
// 					(NewTestProduct("test", "test", "http://test.jp", "", "test_shop", 1111)),
// 					(NewTestProduct("test1", "test1", "http://test.jp", "", "test_shop", 1111)),
// 					(NewTestProduct("test2", "test2", "http://test.jp", "", "test_shop", 1111)),
// 				},
// 				targetProducts: Products{},
// 			},
// 			want: want{
// 				mergedProducts: Products{
// 					(NewTestProduct("test", "test", "http://test.jp", "", "test_shop", 1111)),
// 					(NewTestProduct("test1", "test1", "http://test.jp", "", "test_shop", 1111)),
// 					(NewTestProduct("test2", "test2", "http://test.jp", "", "test_shop", 1111)),
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			merged := tt.args.mergeProducts.MapProducts(tt.args.targetProducts)

// 			assert.Equal(t, tt.want.mergedProducts, merged)
// 		})
// 	}
// }

// func TestGenerateMessage(t *testing.T) {
// 	f := "ikebe_20220301_120303"
// 	type args struct {
// 		product  *Product
// 		filename string
// 	}
// 	type want struct {
// 		message string
// 	}
// 	tests := []struct {
// 		name  string
// 		args  args
// 		want  want
// 		isErr bool
// 	}{
// 		{
// 			name:  "generate message",
// 			args:  args{product: NewTestProduct("test", "test", "https://test.com", "4444", "test_shop", 6000), filename: f},
// 			want:  want{message: `{"filename":"ikebe_20220301_120303","jan":"4444","cost":6000,"url":"https://test.com"}`},
// 			isErr: false,
// 		},
// 		{
// 			name:  "Jan code isn't Valid",
// 			args:  args{product: NewTestProduct("TEST", "test", "https://test.com", "", "test_shop", 5000), filename: f},
// 			want:  want{message: ""},
// 			isErr: true,
// 		},
// 		{
// 			name:  "Price isn't valid",
// 			args:  args{product: NewTestProduct("TEST", "test", "https://test.com", "", "test_shop", 5000), filename: f},
// 			want:  want{message: ""},
// 			isErr: true,
// 		},
// 		{
// 			name:  "URL isn't valid",
// 			args:  args{product: NewTestProduct("TEST", "test", "https://test.com", "", "test_shop", 5000), filename: f},
// 			want:  want{message: ""},
// 			isErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			message, err := tt.args.product.GenerateMessage(tt.args.filename)

// 			assert.Equal(t, tt.want.message, string(message))
// 			if tt.isErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }
