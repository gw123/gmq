package gmq

type Middleware interface {
	GetAttachEventTypes() string
	Handel(event Msg) bool
}
