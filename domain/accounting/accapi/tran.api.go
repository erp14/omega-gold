package accapi

import (
	"net/http"
	"omega/domain/accounting"
	"omega/domain/accounting/accmodel"
	"omega/domain/accounting/accterm"
	"omega/domain/accounting/enum/trantype"
	"omega/domain/service"
	"omega/internal/core"
	"omega/internal/core/corterm"
	"omega/internal/response"
	"omega/pkg/excel"

	"github.com/gin-gonic/gin"
)

// TranAPI for injecting tran service
type TranAPI struct {
	Service service.AccTranServ
	Engine  *core.Engine
}

// ProvideTranAPI for tran is used in wire
func ProvideTranAPI(c service.AccTranServ) TranAPI {
	return TranAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a tran by it's id
func (p *TranAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, accounting.Domain)
	var err error
	var tran accmodel.Tran

	if tran.ID, err = resp.GetRowID(c.Param("tranID"), "E6739759", accterm.Tran); err != nil {
		return
	}

	if tran, err = p.Service.FindByID(tran.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(accounting.ViewTran)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, accterm.Tran).
		JSON(tran)
}

// List of trans
func (p *TranAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, accmodel.TranTable, accounting.Domain)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(accounting.ListTran)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, accterm.Trans).
		JSON(data)
}

// Create tran
func (p *TranAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, accterm.Trans, accounting.Domain)
	var tran, createdTran accmodel.Tran
	var err error

	if err = resp.Bind(&tran, "E6715750", accounting.Domain, accterm.Tran); err != nil {
		return
	}

	tran.Type = trantype.Manual
	tran.CreatedBy = params.UserID

	if createdTran, err = p.Service.Transfer(tran); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(accounting.CreateTran, tran)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, accterm.Tran).
		JSON(createdTran)
}

// Update tran
// func (p *TranAPI) Update(c *gin.Context) {
// 	resp := response.New(p.Engine, c, accounting.Domain)
// 	var err error

// 	var tran, tranBefore, tranUpdated accmodel.Tran

// 	if tran.ID, err = resp.GetRowID(c.Param("tranID"), "E6723735", accterm.Tran); err != nil {
// 		return
// 	}

// 	if err = resp.Bind(&tran, "E6730030", accounting.Domain, accterm.Tran); err != nil {
// 		return
// 	}

// 	if tranBefore, err = p.Service.FindByID(tran.ID); err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	if tranUpdated, err = p.Service.Save(tran); err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	resp.Record(accounting.UpdateTran, tranBefore, tran)
// 	resp.Status(http.StatusOK).
// 		MessageT(corterm.VUpdatedSuccessfully, accterm.Tran).
// 		JSON(tranUpdated)
// }

// Delete tran
// func (p *TranAPI) Delete(c *gin.Context) {
// 	resp := response.New(p.Engine, c, accounting.Domain)
// 	var err error
// 	var tran accmodel.Tran

// 	if tran.ID, err = resp.GetRowID(c.Param("tranID"), "E6744186", accterm.Tran); err != nil {
// 		return
// 	}

// 	if tran, err = p.Service.Delete(tran.ID); err != nil {
// 		resp.Error(err).JSON()
// 		return
// 	}

// 	resp.Record(accounting.DeleteTran, tran)
// 	resp.Status(http.StatusOK).
// 		MessageT(corterm.VDeletedSuccessfully, accterm.Tran).
// 		JSON()
// }

// Excel generate excel files based on search
func (p *TranAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, accterm.Trans, accounting.Domain)

	trans, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("tran")
	ex.AddSheet("Trans").
		AddSheet("Summary").
		Active("Trans").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "B", 15.3).
		SetColWidth("C", "C", 80).
		SetColWidth("D", "E", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Trans").
		WriteHeader("ID", "Name", "Resources", "Description", "Updated At").
		SetSheetFields("ID", "Name", "Resources", "Description", "UpdatedAt").
		WriteData(trans).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(accounting.ExcelTran)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
