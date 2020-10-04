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

// StockAPI for injecting stock service
type StockAPI struct {
	Service service.AccStockServ
	Engine  *core.Engine
}

// ProvideStockAPI for stock is used in wire
func ProvideStockAPI(c service.AccStockServ) StockAPI {
	return StockAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a stock by it's id
func (p *StockAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, accounting.Domain)
	var err error
	var stock accmodel.Stock

	if stock.ID, err = resp.GetRowID(c.Param("stockID"), "E6763178", accterm.Stock); err != nil {
		return
	}

	if stock, err = p.Service.FindByID(stock.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(accounting.ViewStock)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, accterm.Stock).
		JSON(stock)
}

// List of stocks
func (p *StockAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, accmodel.StockTable, accounting.Domain)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(accounting.ListStock)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, accterm.Stocks).
		JSON(data)
}

// Create stock
func (p *StockAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, accounting.Domain)
	var stock, createdStock accmodel.Stock
	var err error

	if err = resp.Bind(&stock, "E6771922", accounting.Domain, accterm.Stock); err != nil {
		return
	}

	if createdStock, err = p.Service.Create(stock); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(accounting.CreateStock, stock)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, accterm.Stock).
		JSON(createdStock)
}

// Update stock
func (p *StockAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, accounting.Domain)
	var err error

	var stock, stockBefore, stockUpdated accmodel.Stock

	if stock.ID, err = resp.GetRowID(c.Param("stockID"), "E6719783", accterm.Stock); err != nil {
		return
	}

	if err = resp.Bind(&stock, "E6781778", accounting.Domain, accterm.Stock); err != nil {
		return
	}

	if stockBefore, err = p.Service.FindByID(stock.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	if stockUpdated, err = p.Service.Save(stock); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(accounting.UpdateStock, stockBefore, stock)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, accterm.Stock).
		JSON(stockUpdated)
}

// Delete stock
func (p *StockAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, accounting.Domain)
	var err error
	var stock accmodel.Stock

	if stock.ID, err = resp.GetRowID(c.Param("stockID"), "E6769002", accterm.Stock); err != nil {
		return
	}

	if stock, err = p.Service.Delete(stock.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(accounting.DeleteStock, stock)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, accterm.Stock).
		JSON()
}

// Excel generate excel files based on search
func (p *StockAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, accterm.Stocks, accounting.Domain)

	stocks, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("stock")
	ex.AddSheet("Stocks").
		AddSheet("Summary").
		Active("Stocks").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "B", 15.3).
		SetColWidth("C", "C", 80).
		SetColWidth("D", "E", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Stocks").
		WriteHeader("ID", "Name", "Resources", "Description", "Updated At").
		SetSheetFields("ID", "Name", "Resources", "Description", "UpdatedAt").
		WriteData(stocks).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(accounting.ExcelStock)

	c.Header("Content-Description", "File Stocksfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
