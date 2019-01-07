package interfaces

type Logger interface {
	Info(category string, content string)
	Warning(category string, content string)
	Error(category string, content string)
	Debug(category string, content string)
}
