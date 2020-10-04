package accounting

import "omega/internal/types"

// list of resources for accounting domain
const (
	Domain string = "accounting"

	TranRead  types.Resource = "tran:read"
	TranWrite types.Resource = "tran:write"
	TranExcel types.Resource = "tran:excel"
)
