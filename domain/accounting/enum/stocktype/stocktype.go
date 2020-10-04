package stocktype

import (
	"omega/internal/types"
	"strings"
)

const (
	CashCurrency    types.Enum = "cash_currency"
	VirtualCurrency types.Enum = "virtual_currency"
	Gold            types.Enum = "gold"
)

var List = []types.Enum{
	CashCurrency,
	VirtualCurrency,
	Gold,
}

// Join make a string for showing in the api
func Join() string {
	var strArr []string

	for _, v := range List {
		strArr = append(strArr, string(v))
	}

	return strings.Join(strArr, ", ")
}
