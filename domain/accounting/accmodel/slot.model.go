package accmodel

import (
	"omega/internal/core/coract"
	"omega/internal/types"
)

// SlotTable is used inside the repo layer
const (
	SlotTable = "acc_slots"
)

// Slot model
type Slot struct {
	types.GormCol
	AccountID   types.RowID `gorm:"not null" json:"account_id,omitempty"`
	TranID      types.RowID `gorm:"not null" json:"tran_id,omitempty"`
	StockID     types.RowID `gorm:"not null" json:"stock_id,omitempty"`
	Debit       float64     `json:"debit"`
	Credit      float64     `json:"credit"`
	Balance     float64     `json:"balance"`
	Description string      `json:"description,omitempty"`
}

// Validate check the type of fields
func (p *Slot) Validate(act coract.Action) (err error) {

	switch {
	case act == coract.Save:

		// if len(p.Name) < 3 {
		// 	err = limberr.AddInvalidParam(err, "name",
		// 		corerr.MinimumAcceptedCharacterForVisV,
		// 		dict.R(corterm.Name), 3)
		// }

		// if len(p.Name) > 255 {
		// 	err = limberr.AddInvalidParam(err, "name",
		// 		corerr.MaximumAcceptedCharacterForVisV,
		// 		dict.R(corterm.Name), 255)
		// }

		// if p.Symbol == "" {
		// 	err = limberr.AddInvalidParam(err, "symbol",
		// 		corerr.VisRequired, dict.R(accterm.Symbol))
		// }

		// if p.Code == "" {
		// 	err = limberr.AddInvalidParam(err, "code",
		// 		corerr.VisRequired, dict.R(corterm.Code))
		// }

		// if len(p.Description) > 255 {
		// 	err = limberr.AddInvalidParam(err, "description",
		// 		corerr.MaximumAcceptedCharacterForVisV,
		// 		dict.R(corterm.Description), 255)
		// }
		// if ok, _ := helper.Includes(slottype.List, p.Type); !ok {
		// 	return limberr.AddInvalidParam(err, "type",
		// 		corerr.AcceptedValueForVareV, dict.R(corterm.Type),
		// 		slottype.Join())
		// }
	}

	return err
}
