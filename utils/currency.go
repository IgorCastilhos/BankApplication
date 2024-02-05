package utils

const (
	BRL = "BRL"
	USD = "USD"
	EUR = "EUR"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case BRL, USD, EUR:
		return true
	}
	return false
}
