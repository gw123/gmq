package gmqcore

import (
	"github.com/gw123/gmq"
)

type ErrorManager struct {
	app gmq.App
}

func NewErrorManager(app gmq.App) *ErrorManager {
	this := new(ErrorManager)
	this.app = app
	return this
}

