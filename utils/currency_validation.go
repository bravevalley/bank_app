package utils

const (
	USD = "USD"
	NGN = "NGN"
	YEN = "YEN"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case USD, NGN, YEN:
		return true
	}

	return false
}
