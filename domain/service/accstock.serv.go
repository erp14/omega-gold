package service

import (
	"fmt"
	"omega/domain/accounting/accmodel"
	"omega/domain/accounting/accrepo"
	"omega/internal/core"
	"omega/internal/core/coract"
	"omega/internal/core/corerr"
	"omega/internal/param"
	"omega/internal/types"
	"omega/pkg/glog"
)

// AccStockServ for injecting auth accrepo
type AccStockServ struct {
	Repo   accrepo.StockRepo
	Engine *core.Engine
}

// ProvideAccStockService for stock is used in wire
func ProvideAccStockService(p accrepo.StockRepo) AccStockServ {
	return AccStockServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting stock by it's id
func (p *AccStockServ) FindByID(id types.RowID) (stock accmodel.Stock, err error) {
	if stock, err = p.Repo.FindByID(id); err != nil {
		err = corerr.Tick(err, "E6726351", "can't fetch the stock", id)
		return
	}

	return
}

// List of stocks, it support pagination and search and return back count
func (p *AccStockServ) List(params param.Param) (stocks []accmodel.Stock,
	count uint64, err error) {

	if stocks, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in stocks list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in stocks count")
	}

	return
}

// Create a stock
func (p *AccStockServ) Create(stock accmodel.Stock) (createdStock accmodel.Stock, err error) {

	if err = stock.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E6790309", "validation failed in creating the stock", stock)
		return
	}

	if createdStock, err = p.Repo.Create(stock); err != nil {
		err = corerr.Tick(err, "E6781681", "stock not created", stock)
		return
	}

	return
}

// Save a stock, if it is exist update it, if not create it
func (p *AccStockServ) Save(stock accmodel.Stock) (savedStock accmodel.Stock, err error) {
	if err = stock.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E6787201", corerr.ValidationFailed, stock)
		return
	}

	if savedStock, err = p.Repo.Save(stock); err != nil {
		err = corerr.Tick(err, "E6735605", "stock not saved")
		return
	}

	return
}

// Delete stock, it is soft delete
func (p *AccStockServ) Delete(stockID types.RowID) (stock accmodel.Stock, err error) {
	if stock, err = p.FindByID(stockID); err != nil {
		err = corerr.Tick(err, "E6759744", "stock not found for deleting")
		return
	}

	if err = p.Repo.Delete(stock); err != nil {
		err = corerr.Tick(err, "E6723097", "stock not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *AccStockServ) Excel(params param.Param) (stocks []accmodel.Stock, err error) {
	params.Limit = p.Engine.Envs.ToUint64(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", accmodel.StockTable)

	if stocks, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E6790567", "cant generate the excel list for stocks")
		return
	}

	return
}
