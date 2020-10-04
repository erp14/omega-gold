package accmodel

import (
	"omega/domain/accounting/enum/trantype"
	"omega/internal/core/coract"
	"omega/internal/core/corerr"
	"omega/internal/core/corterm"
	"omega/internal/types"
	"omega/pkg/dict"
	"omega/pkg/helper"
	"omega/pkg/limberr"
)

// TranTable is used inside the repo layer
const (
	TranTable = "acc_trans"
)

// Tran model
type Tran struct {
	types.GormCol
	PioneerID   types.RowID `json:"pioneer_id,omitempty"`
	FollowerID  types.RowID `json:"follower_id,omitempty"`
	Hash        string      `json:"hash,omitempty"`
	Type        types.Enum  `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	TradeID     types.RowID `json:"trade_id,omitempty"`
}

// Validate check the type of fields
func (p *Tran) Validate(act coract.Action) (err error) {

	// switch act {
	// case coract.Save:

	// 	if len(p.Name) < 5 {
	// 		err = limberr.AddInvalidParam(err, "name",
	// 			corerr.MinimumAcceptedCharacterForVisV,
	// 			dict.R(corterm.Name), 5)
	// 	}

	// 	if len(p.Name) > 255 {
	// 		err = limberr.AddInvalidParam(err, "name",
	// 			corerr.MaximumAcceptedCharacterForVisV,
	// 			dict.R(corterm.Name), 255)
	// 	}

	// 	if p.Resources == "" {
	// 		err = limberr.AddInvalidParam(err, "resources",
	// 			corerr.VisRequired, dict.R(corterm.Resources))
	// 	}

	// 	if len(p.Description) > 255 {
	// 		err = limberr.AddInvalidParam(err, "description",
	// 			corerr.MaximumAcceptedCharacterForVisV,
	// 			dict.R(corterm.Description), 255)
	// 	}
	// }

	// TODO: it should be checked after API has been created
	if ok, _ := helper.Includes(trantype.List, p.Type); !ok {
		// var str []string
		// for _, v := range dict.Langs {
		// 	str = append(str, string(v))
		// }
		return limberr.AddInvalidParam(err, "type",
			corerr.AcceptedValueForVareV, dict.R(corterm.Type),
			trantype.Join())
	}

	return err
}
