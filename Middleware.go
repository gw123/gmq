package gmq2

type Middleware interface {
	GetAttachEventTypes() string
	Handel(event Msg) bool
}
