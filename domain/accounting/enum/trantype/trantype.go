package trantype

import (
	"omega/internal/types"
	"strings"
)

const (
	Trade   types.Enum = "trade"
	Asset   types.Enum = "asset"
	Expense types.Enum = "expense"
	Manual  types.Enum = "manual"
)

var List = []types.Enum{
	Trade,
	Asset,
	Expense,
	Manual,
}

// Join make a string for showing in the api
func Join() string {
	var strArr []string

	for _, v := range List {
		strArr = append(strArr, string(v))
	}

	return strings.Join(strArr, ", ")
}
