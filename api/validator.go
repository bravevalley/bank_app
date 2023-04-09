package api

import (
	"github.com/dassyareg/bank_app/utils"
	"github.com/go-playground/validator/v10"
)

var validateCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if curr, ok := fl.Field().Interface().(string); ok {
		return utils.IsValidCurrency(curr)
	}
	return false
}