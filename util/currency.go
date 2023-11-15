package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	MYR = "MYR"
	SGD = "SGD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, MYR, SGD:
		return true
	}
	return false
}
