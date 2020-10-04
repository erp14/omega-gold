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

// AccSlotServ for injecting auth accrepo
type AccSlotServ struct {
	Repo   accrepo.SlotRepo
	Engine *core.Engine
}

// ProvideAccSlotService for slot is used in wire
func ProvideAccSlotService(p accrepo.SlotRepo) AccSlotServ {
	return AccSlotServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting slot by it's id
func (p *AccSlotServ) FindByID(id types.RowID) (slot accmodel.Slot, err error) {
	if slot, err = p.Repo.FindByID(id); err != nil {
		err = corerr.Tick(err, "E6741209", "can't fetch the slot", id)
		return
	}

	return
}

// List of slots, it support pagination and search and return back count
func (p *AccSlotServ) List(params param.Param) (slots []accmodel.Slot,
	count uint64, err error) {

	if slots, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in slots list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in slots count")
	}

	return
}

// Create a slot
func (p *AccSlotServ) Create(slot accmodel.Slot) (createdSlot accmodel.Slot, err error) {

	if err = slot.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E6799673", "validation failed in creating the slot", slot)
		return
	}

	if createdSlot, err = p.Repo.Create(slot); err != nil {
		err = corerr.Tick(err, "E6734399", "slot not created", slot)
		return
	}

	return
}

// Save a slot, if it is exist update it, if not create it
func (p *AccSlotServ) Save(slot accmodel.Slot) (savedSlot accmodel.Slot, err error) {
	if err = slot.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E6795506", corerr.ValidationFailed, slot)
		return
	}

	if savedSlot, err = p.Repo.Save(slot); err != nil {
		err = corerr.Tick(err, "E6793602", "slot not saved")
		return
	}

	return
}

// Delete slot, it is soft delete
func (p *AccSlotServ) Delete(slotID types.RowID) (slot accmodel.Slot, err error) {
	if slot, err = p.FindByID(slotID); err != nil {
		err = corerr.Tick(err, "E6782480", "slot not found for deleting")
		return
	}

	if err = p.Repo.Delete(slot); err != nil {
		err = corerr.Tick(err, "E6715070", "slot not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *AccSlotServ) Excel(params param.Param) (slots []accmodel.Slot, err error) {
	params.Limit = p.Engine.Envs.ToUint64(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", accmodel.SlotTable)

	if slots, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E6750180", "cant generate the excel list for slots")
		return
	}

	return
}
