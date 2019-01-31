package interfaces

type Logger interface {
	Info(category string, format string, a ...interface{})
	Warning(category string, format string, a ...interface{})
	Error(category string, format string, a ...interface{})
	Debug(category string, format string, a ...interface{})
}
