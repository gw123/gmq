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
	db.AutoMigrate(&models.Service{})
	db.AutoMigrate(&models.Client{})
	db.AutoMigrate(&models.ClientServer{})

	return nil
}
