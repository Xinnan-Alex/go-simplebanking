package api

import (
	"github.com/Xinnan-Alex/simplebank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	currency, ok := fieldLevel.Field().Interface().(string)
	if ok {
		//check curency is supported
		return util.IsSupportedCurrency(currency)
	}
	return false
}
