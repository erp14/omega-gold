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

// StockRepo for injecting engine
type StockRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideStockRepo is used in wire and initiate the Cols
func ProvideStockRepo(engine *core.Engine) StockRepo {
	return StockRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(accmodel.Stock{}), accmodel.StockTable),
	}
}

// FindByID finds the stock via its id
func (p *StockRepo) FindByID(id types.RowID) (stock accmodel.Stock, err error) {
	err = p.Engine.DB.Table(accmodel.StockTable).First(&stock, id.ToUint64()).Error

	stock.ID = id
	err = p.dbError(err, "E6718084", stock, corterm.List)

	return
}

// List returns an array of stocks
func (p *StockRepo) List(params param.Param) (stocks []accmodel.Stock, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E6764143").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E6745113").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.DB.Table(accmodel.StockTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&stocks).Error

	err = p.dbError(err, "E6788036", accmodel.Stock{}, corterm.List)

	return
}

// Count of stocks, mainly calls with List
func (p *StockRepo) Count(params param.Param) (count uint64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E6780723").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.DB.Table(accmodel.StockTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E6744054", accmodel.Stock{}, corterm.List)
	return
}

// Save the stock, in case it is not exist create it
func (p *StockRepo) Save(stock accmodel.Stock) (u accmodel.Stock, err error) {
	if err = p.Engine.DB.Table(accmodel.StockTable).Save(&stock).Error; err != nil {
		err = p.dbError(err, "E6717660", stock, corterm.Updated)
	}

	p.Engine.DB.Table(accmodel.StockTable).Where("id = ?", stock.ID).Find(&u)
	return
}

// Create a stock
func (p *StockRepo) Create(stock accmodel.Stock) (u accmodel.Stock, err error) {
	if err = p.Engine.DB.Table(accmodel.StockTable).Create(&stock).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E6714072", stock, corterm.Created)
	}
	return
}

// Delete the stock
func (p *StockRepo) Delete(stock accmodel.Stock) (err error) {
	if err = p.Engine.DB.Table(accmodel.StockTable).Unscoped().Delete(&stock).Error; err != nil {
		err = p.dbError(err, "E6719633", stock, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper database error
func (p *StockRepo) dbError(err error, code string, stock accmodel.Stock, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, stock.ID, accterm.Stocks)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(basterm.Users),
				dict.R(accterm.Stock), dict.R(action)).
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
