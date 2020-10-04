package accapi

import (
	"net/http"
	"omega/domain/accounting"
	"omega/domain/accounting/accmodel"
	"omega/domain/accounting/accterm"
	"omega/domain/service"
	"omega/internal/core"
	"omega/internal/core/corterm"
	"omega/internal/response"
	"omega/pkg/excel"

	"github.com/gin-gonic/gin"
)

// SlotAPI for injecting slot service
type SlotAPI struct {
	Service service.AccSlotServ
	Engine  *core.Engine
}

// ProvideSlotAPI for slot is used in wire
func ProvideSlotAPI(c service.AccSlotServ) SlotAPI {
	return SlotAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a slot by it's id
func (p *SlotAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, accounting.Domain)
	var err error
	var slot accmodel.Slot

	if slot.ID, err = resp.GetRowID(c.Param("slotID"), "E6750783", accterm.Slot); err != nil {
		return
	}

	if slot, err = p.Service.FindByID(slot.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(accounting.ViewSlot)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, accterm.Slot).
		JSON(slot)
}

// List of slots
func (p *SlotAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, accmodel.SlotTable, accounting.Domain)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(accounting.ListSlot)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, accterm.Slots).
		JSON(data)
}

// Create slot
func (p *SlotAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, accounting.Domain)
	var slot, createdSlot accmodel.Slot
	var err error

	if err = resp.Bind(&slot, "E6753711", accounting.Domain, accterm.Slot); err != nil {
		return
	}

	if createdSlot, err = p.Service.Create(slot); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(accounting.CreateSlot, slot)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, accterm.Slot).
		JSON(createdSlot)
}

// Update slot
func (p *SlotAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, accounting.Domain)
	var err error

	var slot, slotBefore, slotUpdated accmodel.Slot

	if slot.ID, err = resp.GetRowID(c.Param("slotID"), "E6746218", accterm.Slot); err != nil {
		return
	}

	if err = resp.Bind(&slot, "E6733638", accounting.Domain, accterm.Slot); err != nil {
		return
	}

	if slotBefore, err = p.Service.FindByID(slot.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	if slotUpdated, err = p.Service.Save(slot); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(accounting.UpdateSlot, slotBefore, slot)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, accterm.Slot).
		JSON(slotUpdated)
}

// Delete slot
func (p *SlotAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, accounting.Domain)
	var err error
	var slot accmodel.Slot

	if slot.ID, err = resp.GetRowID(c.Param("slotID"), "E6737806", accterm.Slot); err != nil {
		return
	}

	if slot, err = p.Service.Delete(slot.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(accounting.DeleteSlot, slot)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, accterm.Slot).
		JSON()
}

// Excel generate excel files based on search
func (p *SlotAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, accterm.Slots, accounting.Domain)

	slots, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("slot")
	ex.AddSheet("Slots").
		AddSheet("Summary").
		Active("Slots").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "B", 15.3).
		SetColWidth("C", "C", 80).
		SetColWidth("D", "E", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Slots").
		WriteHeader("ID", "Name", "Resources", "Description", "Updated At").
		SetSheetFields("ID", "Name", "Resources", "Description", "UpdatedAt").
		WriteData(slots).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(accounting.ExcelSlot)

	c.Header("Content-Description", "File Slotsfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
