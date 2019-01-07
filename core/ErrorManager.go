package core

import (
	"github.com/gw123/GMQ/core/interfaces"
)

type ErrorManager struct {
	app interfaces.App
}

func NewErrorManager(app interfaces.App) *ErrorManager {
	this := new(ErrorManager)
	this.app = app
	return this
}
