package tutorial

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Getter[T any] interface {
	Get() T
	Set(T) Getter[T]
}
type MyStruct[T any] struct {
	Val T
}

func (m MyStruct[T]) Get() T {
	return m.Val
}

func (m MyStruct[T]) Set(val T) Getter[T] {
	m.Val = val
	return m
}

func bar[T any]() Getter[T] {
	return MyStruct[T]{}
}

func TestInterface(t *testing.T) {
	mystruct := bar[string]()
	mystruct = mystruct.Set("not found error")
	fmt.Println(mystruct.Get())
}

func TestCond(t *testing.T) {
	type cond struct {
		isManaged *bool
	}
	c := cond{}
	assert.True(t, c.isManaged == nil || *c.isManaged)

	c1 := cond{
		isManaged: Ptr(false),
	}
	assert.False(t, c1.isManaged != nil && *c1.isManaged)
}
