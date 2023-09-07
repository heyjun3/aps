package bomber

import (
	"crawler/test/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduct(t *testing.T) {
	type args struct {
		filename string
	}
	type want struct {
		code string
	}

	tests := []struct {
		name  string
		args  args
		want  want
		isErr bool
	}{{
		name:  "parse product page",
		args:  args{filename: "html/test_product.html"},
		want:  want{code: "4550161577867"},
		isErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CreateHttpResponse(tt.args.filename)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			code, err := BomberParser{}.Product(res.Body)

			assert.Equal(t, tt.want.code, code)
			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
