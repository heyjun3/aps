package tutorial

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceLoop(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	ex := []*int{Ptr(1), Ptr(2), Ptr(3), Ptr(4), Ptr(5)}

	r := make([]*int, 0, len(s))
	for _, i := range s {
		x := i
		r = append(r, &x)
	}

	assert.Equal(t, ex, r)
}
