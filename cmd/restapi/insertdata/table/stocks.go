package table

import (
	"omega/domain/accounting/accmodel"
	"omega/domain/accounting/accrepo"
	"omega/domain/accounting/enum/stocktype"
	"omega/domain/service"
	"omega/internal/core"
	"omega/internal/types"
	"omega/pkg/glog"
)

// InsertStocks for add required stocks
func InsertStocks(engine *core.Engine) {
	engine.DB.Exec("UPDATE bas_stocks SET deleted_at = null WHERE id IN (1,2,3,4)")
	stockRepo := accrepo.ProvideStockRepo(engine)
	stockService := service.ProvideAccStockService(stockRepo)
	stocks := []accmodel.Stock{
		{
			GormCol: types.GormCol{
				ID: 1,
			},
			Name:        "Dollar Cash",
			Type:        stocktype.CashCurrency,
			Symbol:      "$",
			Code:        "USD-Cash",
			Description: "usd cash for getting physically",
		},
		{
			GormCol: types.GormCol{
				ID: 2,
			},
			Name:        "Dollar Virtual",
			Type:        stocktype.VirtualCurrency,
			Symbol:      "$",
			Code:        "USD-V",
			Description: "usd cash for trades",
		},
		{
			GormCol: types.GormCol{
				ID: 3,
			},
			Name:        "IQ Gold-18",
			Type:        stocktype.Gold,
			Symbol:      "G18",
			Code:        "IQ-G18",
			Description: "gold 18",
		},
		{
			GormCol: types.GormCol{
				ID: 4,
			},
			Name:        "IQ Gold-24",
			Type:        stocktype.Gold,
			Symbol:      "G24",
			Code:        "IQ-G24",
			Description: "gold 24",
		},
	}

	for _, v := range stocks {
		if _, err := stockService.Save(v); err != nil {
			glog.Fatal("error in saving stocks", err)
		}
	}

}
