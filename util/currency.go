package util

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	NIS = "NIS"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, NIS:
		return true
	}
	return false
}
