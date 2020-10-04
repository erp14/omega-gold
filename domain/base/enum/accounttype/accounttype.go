package accounttype

import (
	"omega/internal/types"
	"strings"
)

const (
	Eternal  types.Enum = "eternal"
	Asset    types.Enum = "asset"
	Expense  types.Enum = "expense"
	Trader   types.Enum = "trader"
	Provider types.Enum = "provider"
	Cashier  types.Enum = "cashier"
	Fee      types.Enum = "fee"
	Fixer    types.Enum = "fixer"
)

var List = []types.Enum{
	Eternal,
	Asset,
	Expense,
	Trader,
	Provider,
	Cashier,
	Fee,
	Fixer,
}

// Join make a string for showing in the api
func Join() string {
	var strArr []string

	for _, v := range List {
		strArr = append(strArr, string(v))
	}

	return strings.Join(strArr, ", ")
}
