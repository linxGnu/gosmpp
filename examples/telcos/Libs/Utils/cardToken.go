package Utils

import (
	"strings"
)

var internationalCardScheme = []string{"VISA", "MASTERCARD", "AMEX", "JCB"}

// IsATMCard check if card scheme is atm card or not
func IsATMCard(cardScheme string) bool {
	cardScheme = strings.TrimSpace(cardScheme)

	for i := range internationalCardScheme {
		if cardScheme == internationalCardScheme[i] {
			return false
		}
	}

	return true
}
