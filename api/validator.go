package api

import "github.com/go-playground/validator/v10"

// Currency constants
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	VND = "VND"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)

	if ok {
		return IsSupportCurrency(currency)
	}

	return false
}

func IsSupportCurrency(currency string) bool {
	switch currency {
	case CAD, USD, EUR, VND:
		return true
	}

	return false
}
