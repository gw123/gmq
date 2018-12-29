package interfaces

type Middleware interface {
	GetAttachEventTypes() string
	Handel(event Event) bool
}
