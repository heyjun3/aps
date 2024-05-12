package tutorial

import (
	"fmt"
	"sort"
	"testing"
	"time"

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

func TestSortDateSlice(t *testing.T) {
	dates := []string{
		"2024-05-19",
		"2024-05-11",
		"2024-05-12",
		"2024-05-14",
	}

	const format string = "2006-01-02"

	sort.Slice(dates, func(i, j int) bool {
		ti, err := time.Parse(format, dates[i])
		if err != nil {
			panic(err)
		}
		tj, err := time.Parse(format, dates[j])
		if err != nil {
			panic(err)
		}
		return ti.Before(tj)
	})
	fmt.Println(dates[1:2])
}
