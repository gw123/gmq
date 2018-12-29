package app

import (
	"github.com/gw123/GMQ/interfaces"
)

type ErrorManager struct {
	app interfaces.App
}

func NewErrorManager(app interfaces.App) *ErrorManager {
	this := new(ErrorManager)
	this.app = app
	return this
}
