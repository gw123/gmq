package bootstarp

import (
	db_models "github.com/gw123/GMQ/common/models"
	"github.com/gw123/GMQ/core/interfaces"
)

func AutoMigrate(app interfaces.App) error {
	db, err := app.GetDefaultDb()
	if err != nil {
		return err
	}
	app.Info("App", "迁移数据库")
	//db.AutoMigrate(&models.Task{})
	//db.AutoMigrate(&models.TaskDetail{})
	//db.AutoMigrate(&models.Client{})
	//db.AutoMigrate(&models.ClientTask{})
	//db.AutoMigrate(&models.PingLog{})
	//db.AutoMigrate(&models.Comment{})

	c := db_models.Comment{
		Ip: "127.0.0.1",
	}
	result := db.Save(&c).Error
	if result != nil {
		app.Error("App", result.Error())
	}
	app.Info("App", "迁移完成")
	return nil
}
