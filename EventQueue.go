package gmq2

type EventQueue interface {
	Push(event Msg) error
	Pop() (Msg, error)
}
