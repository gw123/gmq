package interfaces

type ModuleLogger interface {
	Info(format string, a ...interface{})
	Warning(format string, a ...interface{})
	Error(format string, a ...interface{})
	Debug(format string, a ...interface{})
}
