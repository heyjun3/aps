package scrape

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeToStr(t *testing.T) {
	t.Run("format time to str", func(t *testing.T) {
		d := time.Date(2023, 2, 9, 22, 59, 0, 0, time.Local)

		s := TimeToStr(d)

		assert.Equal(t, "20230209_225900", s)
	})
}

func TestPullOutPrice(t *testing.T) {
	type args struct {
		str string
	}
	type want struct {
		num int64
	}
	tests := []struct {
		name  string
		args  args
		want  want
		isErr bool
	}{
		{
			name:  "pull out price",
			args:  args{str: " 199,800円"},
			want:  want{num: 199800},
			isErr: false,
		},
		{
			name:  "pull out price not digits",
			args:  args{str: "aaa  fdsagfda"},
			want:  want{num: 0},
			isErr: true,
		},
		{
			name:  "blank string",
			args:  args{str: ""},
			want:  want{num: 0},
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			num, err := PullOutNumber(tt.args.str)

			assert.Equal(t, tt.want.num, num)
			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
