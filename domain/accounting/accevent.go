package accounting

import "omega/internal/types"

// types for accounting domain
const (
	CreateTran types.Event = "tran-create"
	UpdateTran types.Event = "tran-update"
	DeleteTran types.Event = "tran-delete"
	ListTran   types.Event = "tran-list"
	ViewTran   types.Event = "tran-view"
	ExcelTran  types.Event = "tran-excel"
)
