package tutorial

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	t.Run("validate required filed", func(t *testing.T) {
		type Shop struct {
			Id   string `validate:"required"`
			Name string `validate:"required"`
		}
		shop := Shop{}
		validate := validator.New()

		err := validate.Struct(shop)

		assert.Error(t, err)
	})

	t.Run("validate required field with rules", func(t *testing.T) {
		type Shop struct {
			Id   string
			Name string
		}
		shop := Shop{}
		rules := map[string]string{"Id": "required", "Name": "required"}
		validate := validator.New()
		validate.RegisterStructValidationMapRules(rules, Shop{})

		err := validate.Struct(shop)

		assert.Error(t, err)
	})

	t.Run("validate nest struct", func(t *testing.T) {
		type Shop struct {
			Id   string `validate:"required"`
			Name string
		}
		type Shops struct {
			Id   string
			Shop []Shop `validate:"dive"`
		}
		shops := Shops{Shop: []Shop{{Name: "test"}}}
		validate := validator.New()

		err := validate.Struct(shops)

		assert.Error(t, err)
	})
}
