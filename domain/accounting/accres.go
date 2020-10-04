package accounting

import "omega/internal/types"

// list of resources for accounting domain
const (
	Domain string = "accounting"

	StockRead  types.Resource = "stock:read"
	StockWrite types.Resource = "stock:write"
	StockExcel types.Resource = "stock:excel"

	TranRead  types.Resource = "tran:read"
	TranWrite types.Resource = "tran:write"
	TranExcel types.Resource = "tran:excel"
)
