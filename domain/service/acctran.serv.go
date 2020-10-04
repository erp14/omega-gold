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

// AccTranServ for injecting auth accrepo
type AccTranServ struct {
	Repo   accrepo.TranRepo
	Engine *core.Engine
}

// ProvideAccTranService for tran is used in wire
func ProvideAccTranService(p accrepo.TranRepo) AccTranServ {
	return AccTranServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting tran by it's id
func (p *AccTranServ) FindByID(id types.RowID) (tran accmodel.Tran, err error) {
	if tran, err = p.Repo.FindByID(id); err != nil {
		err = corerr.Tick(err, "E6770185", "can't fetch the tran", id)
		return
	}

	return
}

// List of trans, it support pagination and search and return back count
func (p *AccTranServ) List(params param.Param) (trans []accmodel.Tran,
	count uint64, err error) {

	if trans, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in trans list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in trans count")
	}

	return
}

// Create a tran
func (p *AccTranServ) Create(tran accmodel.Tran) (createdTran accmodel.Tran, err error) {

	if err = tran.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E6771623", "validation failed in creating the tran", tran)
		return
	}

	if createdTran, err = p.Repo.Create(tran); err != nil {
		err = corerr.Tick(err, "E6759336", "tran not created", tran)
		return
	}

	return
}

// Save a tran, if it is exist update it, if not create it
func (p *AccTranServ) Save(tran accmodel.Tran) (savedTran accmodel.Tran, err error) {
	if err = tran.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E6751314", corerr.ValidationFailed, tran)
		return
	}

	if savedTran, err = p.Repo.Save(tran); err != nil {
		err = corerr.Tick(err, "E6737536", "tran not saved")
		return
	}

	return
}

// Delete tran, it is soft delete
func (p *AccTranServ) Delete(tranID types.RowID) (tran accmodel.Tran, err error) {
	if tran, err = p.FindByID(tranID); err != nil {
		err = corerr.Tick(err, "E6753169", "tran not found for deleting")
		return
	}

	if err = p.Repo.Delete(tran); err != nil {
		err = corerr.Tick(err, "E6795168", "tran not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *AccTranServ) Excel(params param.Param) (trans []accmodel.Tran, err error) {
	params.Limit = p.Engine.Envs.ToUint64(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", accmodel.TranTable)

	if trans, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E6742933", "cant generate the excel list for trans")
		return
	}

	return
}
