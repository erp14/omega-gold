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
)

// TranRepo for injecting engine
type TranRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideTranRepo is used in wire and initiate the Cols
func ProvideTranRepo(engine *core.Engine) TranRepo {
	return TranRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(accmodel.Tran{}), accmodel.TranTable),
	}
}

// FindByID finds the tran via its id
func (p *TranRepo) FindByID(id types.RowID) (tran accmodel.Tran, err error) {
	err = p.Engine.DB.Table(accmodel.TranTable).First(&tran, id.ToUint64()).Error

	tran.ID = id
	err = p.dbError(err, "E6744242", tran, corterm.List)

	return
}

// List returns an array of trans
func (p *TranRepo) List(params param.Param) (trans []accmodel.Tran, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E6765798").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.DB.Table(accmodel.TranTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&trans).Error

	err = p.dbError(err, "E6790433", accmodel.Tran{}, corterm.List)

	return
}

// Count of trans, mainly calls with List
func (p *TranRepo) Count(params param.Param) (count uint64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E6795694").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.DB.Table(accmodel.TranTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E6780301", accmodel.Tran{}, corterm.List)
	return
}

// Save the tran, in case it is not exist create it
// func (p *TranRepo) Save(tran accmodel.Tran) (u accmodel.Tran, err error) {
// 	if err = p.Engine.DB.Table(accmodel.TranTable).Save(&tran).Error; err != nil {
// 		err = p.dbError(err, "E6716310", tran, corterm.Updated)
// 	}

// 	p.Engine.DB.Table(accmodel.TranTable).Where("id = ?", tran.ID).Find(&u)
// 	return
// }

// Create a tran
func (p *TranRepo) Create(tran accmodel.Tran) (u accmodel.Tran, err error) {
	if err = p.Engine.DB.Table(accmodel.TranTable).Create(&tran).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E6751309", tran, corterm.Created)
	}
	return
}

// Delete the tran
// func (p *TranRepo) Delete(tran accmodel.Tran) (err error) {
// 	if err = p.Engine.DB.Table(accmodel.TranTable).Unscoped().Delete(&tran).Error; err != nil {
// 		err = p.dbError(err, "E6775100", tran, corterm.Deleted)
// 	}
// 	return
// }

// dbError is an internal method for generate proper database error
func (p *TranRepo) dbError(err error, code string, tran accmodel.Tran, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, tran.ID, accterm.Trans)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(basterm.Users),
				dict.R(accterm.Tran), dict.R(action)).
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
