package tutorial

import (
	"fmt"
	"golang.org/x/exp/slices"
	"reflect"
	"testing"
	// "github.com/stretchr/testify/assert"
)

func Ptr[T any](v T) *T {
	return &v
}

func TestIterateFieldsOfStruct(t *testing.T) {
	type Foo struct {
		Foo *string `bun:"foo"`
		Bar *string
	}

	t.Run("iterate fileds of Foo struct", func(t *testing.T) {
		foo := Foo{Foo: Ptr("foo"), Bar: Ptr("bar")}

		v := reflect.ValueOf(foo)
		types := v.Type()
		fmt.Println(v)

		for i := 0; i < v.NumField(); i++ {
			t := (types.Field(i))
			fmt.Println(t.Name)
			fmt.Println(t.Tag.Get("bun"))
			fmt.Println(v.Field(i).Interface() == nil)
		}
	})

	t.Run("contains nil in slice", func(t *testing.T) {
		s := []interface{}{nil, nil}
		fmt.Println(slices.Contains(s, nil))
	})
}
