package accounting

import "omega/internal/types"

// types for accounting domain
const (
	CreateStock types.Event = "stock-create"
	UpdateStock types.Event = "stock-update"
	DeleteStock types.Event = "stock-delete"
	ListStock   types.Event = "stock-list"
	ViewStock   types.Event = "stock-view"
	ExcelStock  types.Event = "stock-excel"

	CreateTran types.Event = "tran-create"
	UpdateTran types.Event = "tran-update"
	DeleteTran types.Event = "tran-delete"
	ListTran   types.Event = "tran-list"
	ViewTran   types.Event = "tran-view"
	ExcelTran  types.Event = "tran-excel"

	CreateSlot types.Event = "slot-create"
	UpdateSlot types.Event = "slot-update"
	DeleteSlot types.Event = "slot-delete"
	ListSlot   types.Event = "slot-list"
	ViewSlot   types.Event = "slot-view"
	ExcelSlot  types.Event = "slot-excel"
)
