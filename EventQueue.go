package gmq

type EventQueue interface {
	Push(event Msg) error
	Pop() (Msg, error)
}
