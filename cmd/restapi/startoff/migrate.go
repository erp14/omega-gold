package startoff

import (
	"fmt"
	"omega/domain/accounting/accmodel"
	"omega/domain/base/basmodel"
	"omega/internal/core"
)

// Migrate the database for creating tables
func Migrate(engine *core.Engine) {
	// Base Domain
	engine.DB.Table(basmodel.SettingTable).AutoMigrate(&basmodel.Setting{})
	engine.DB.Table(basmodel.RoleTable).AutoMigrate(&basmodel.Role{})
	engine.DB.Table(basmodel.UserTable).AutoMigrate(&basmodel.User{}).
		AddForeignKey("role_id", fmt.Sprintf("%v(id)", basmodel.RoleTable), "RESTRICT", "RESTRICT")
	engine.ActivityDB.Table(basmodel.ActivityTable).AutoMigrate(&basmodel.Activity{})

	engine.DB.Table(basmodel.AccountTable).AutoMigrate(&basmodel.Account{})

	// Accounting Domain
	engine.DB.Table(accmodel.StockTable).AutoMigrate(&accmodel.Stock{})
	engine.DB.Table(accmodel.TranTable).AutoMigrate(&accmodel.Tran{}).
		AddForeignKey("pioneer_id", "bas_accounts(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("follower_id", "bas_accounts(id)", "RESTRICT", "RESTRICT")
	// AddForeignKey("trade_id", "(id)", "RESTRICT", "RESTRICT")

}
