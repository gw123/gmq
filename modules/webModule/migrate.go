package webModule

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/webModule/db_models"
)

func autoMigrate(app interfaces.App) error {
	db, err := app.GetDefaultDb()
	if err != nil {
		return err
	}
	db.AutoMigrate(&db_models.Task{})
	db.AutoMigrate(&db_models.TaskDetail{})
	db.AutoMigrate(&db_models.Client{})
	db.AutoMigrate(&db_models.ClientTask{})
	db.AutoMigrate(&db_models.PingLog{})
	return nil
}
