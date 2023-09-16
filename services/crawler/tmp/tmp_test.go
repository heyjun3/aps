package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStruct(t *testing.T) {
	type Cat struct {
		name string
	}
	type Dog struct {
		name string
	}
	dog := Dog{name: "tom"}
	cat := Cat{name: "tom"}
	// cat2 := Dog{name: cat.name}
	cat2 := Dog(cat)
	assert.Equal(t, cat2, dog)
}
