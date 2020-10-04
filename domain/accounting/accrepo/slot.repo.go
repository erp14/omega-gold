package accrepo

import (
	"omega/domain/accounting/accmodel"
	"omega/domain/accounting/accterm"
	"omega/domain/base/message/basterm"
	"omega/internal/core"
	"omega/internal/core/corerr"
	"omega/internal/core/corterm"
	"omega/internal/core/validator"
	"omega/internal/param"
	"omega/internal/types"
	"omega/pkg/dict"
	"omega/pkg/helper"
	"omega/pkg/limberr"
	"reflect"

	"github.com/jinzhu/gorm"
)

// SlotRepo for injecting engine
type SlotRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideSlotRepo is used in wire and initiate the Cols
func ProvideSlotRepo(engine *core.Engine) SlotRepo {
	return SlotRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(accmodel.Slot{}), accmodel.SlotTable),
	}
}

// FindByID finds the slot via its id
func (p *SlotRepo) FindByID(id types.RowID) (slot accmodel.Slot, err error) {
	err = p.Engine.DB.Table(accmodel.SlotTable).First(&slot, id.ToUint64()).Error

	slot.ID = id
	err = p.dbError(err, "E6739358", slot, corterm.List)

	return
}

// LastSlot return last balance for the account
func (p *SlotRepo) LastSlot(accountID, stockID types.RowID) (slot accmodel.Slot, err error) {
	err = p.Engine.DB.Table(accmodel.SlotTable).
		Where("account_id = ? && stock_id = ?", accountID, stockID).
		Find(&slot).Error

	if gorm.IsRecordNotFoundError(err) {
		err = nil
	}

	slot.AccountID = accountID
	slot.StockID = stockID
	err = p.dbError(err, "E6739358", slot, corterm.List)

	return
}

// List returns an array of slots
func (p *SlotRepo) List(params param.Param) (slots []accmodel.Slot, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E6756148").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E6778533").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.DB.Table(accmodel.SlotTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&slots).Error

	err = p.dbError(err, "E6724947", accmodel.Slot{}, corterm.List)

	return
}

// Count of slots, mainly calls with List
func (p *SlotRepo) Count(params param.Param) (count uint64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E6724242").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.DB.Table(accmodel.SlotTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E6786291", accmodel.Slot{}, corterm.List)
	return
}

// Save the slot, in case it is not exist create it
// func (p *SlotRepo) Save(slot accmodel.Slot) (u accmodel.Slot, err error) {
// 	if err = p.Engine.DB.Table(accmodel.SlotTable).Save(&slot).Error; err != nil {
// 		err = p.dbError(err, "E6760034", slot, corterm.Updated)
// 	}

// 	p.Engine.DB.Table(accmodel.SlotTable).Where("id = ?", slot.ID).Find(&u)
// 	return
// }

// Create a slot
func (p *SlotRepo) Create(slot accmodel.Slot) (u accmodel.Slot, err error) {
	if err = p.Engine.DB.Table(accmodel.SlotTable).Create(&slot).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E6721225", slot, corterm.Created)
	}
	return
}

// Delete the slot
// func (p *SlotRepo) Delete(slot accmodel.Slot) (err error) {
// 	if err = p.Engine.DB.Table(accmodel.SlotTable).Unscoped().Delete(&slot).Error; err != nil {
// 		err = p.dbError(err, "E6727122", slot, corterm.Deleted)
// 	}
// 	return
// }

// dbError is an internal method for generate proper database error
func (p *SlotRepo) dbError(err error, code string, slot accmodel.Slot, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, slot.ID, accterm.Slots)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(basterm.Users),
				dict.R(accterm.Slot), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
