package middlewares

import "github.com/gw123/GMQ/core/interfaces"

type EventAuth struct {
	app interfaces.App
}

func NewEventAuth(app interfaces.App) *EventAuth {
	this := new(EventAuth)
	this.app = app
	return this
}
func (this *EventAuth) Handel(event interfaces.Event) bool {
	//this.app.Debug("eventAuth", event.GetMsgId()+" : "+event.GetEventName())
	return true
}

func (this *EventAuth) GetAttachEventTypes() string {
	return "*"
}
