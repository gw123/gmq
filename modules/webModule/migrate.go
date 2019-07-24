package webModule

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/webModule/models"
)

func autoMigrate(app interfaces.App) error {
	db, err := app.GetDefaultDb()
	if err != nil {
		return err
	}
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.TaskDetail{})
	db.AutoMigrate(&models.Client{})
	db.AutoMigrate(&models.ClientTask{})
	db.AutoMigrate(&models.PingLog{})
	return nil
}
