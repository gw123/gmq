package middlewares

import "github.com/gw123/GMQ/interfaces"

type EventFormat struct {
	app interfaces.App
}

func NewEventFormat(app interfaces.App) *EventFormat {
	this := new(EventFormat)
	this.app = app
	return this
}

func (this *EventFormat) Handel(event interfaces.Event) bool {
	//this.app.Debug("eventFormat", event.GetMsgId()+":"+event.GetEventName())
	return true
}

func (this *EventFormat) GetAttachEventTypes() string {
	return "*"
}
