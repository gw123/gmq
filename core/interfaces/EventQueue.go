package interfaces

type EventQueue interface {
	Push(event Event) error
	Pop() (Event, error)
}
