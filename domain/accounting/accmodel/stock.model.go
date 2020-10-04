package accmodel

import (
	"omega/domain/accounting/accterm"
	"omega/domain/accounting/enum/stocktype"
	"omega/internal/core/coract"
	"omega/internal/core/corerr"
	"omega/internal/core/corterm"
	"omega/internal/types"
	"omega/pkg/dict"
	"omega/pkg/helper"
	"omega/pkg/limberr"
)

// StockTable is used inside the repo layer
const (
	StockTable = "acc_stocks"
)

// Stock model
type Stock struct {
	types.GormCol
	Name        string     `gorm:"not null;unique" json:"name,omitempty"`
	Type        types.Enum `gorm:"not null" json:"type,omitempty"`
	Symbol      string     `gorm:"not null" json:"symbol,omitempty"`
	Code        string     `gorm:"not null;unique" json:"code,omitempty"`
	Description string     `json:"description,omitempty"`
}

// Validate check the type of fields
func (p *Stock) Validate(act coract.Action) (err error) {

	switch {
	case act == coract.Save:

		if len(p.Name) < 3 {
			err = limberr.AddInvalidParam(err, "name",
				corerr.MinimumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 3)
		}

		if len(p.Name) > 255 {
			err = limberr.AddInvalidParam(err, "name",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 255)
		}

		if p.Symbol == "" {
			err = limberr.AddInvalidParam(err, "symbol",
				corerr.VisRequired, dict.R(accterm.Symbol))
		}

		if p.Code == "" {
			err = limberr.AddInvalidParam(err, "code",
				corerr.VisRequired, dict.R(corterm.Code))
		}

		if len(p.Description) > 255 {
			err = limberr.AddInvalidParam(err, "description",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Description), 255)
		}
		if ok, _ := helper.Includes(stocktype.List, p.Type); !ok {
			return limberr.AddInvalidParam(err, "type",
				corerr.AcceptedValueForVareV, dict.R(corterm.Type),
				stocktype.Join())
		}
	}

	return err
}
