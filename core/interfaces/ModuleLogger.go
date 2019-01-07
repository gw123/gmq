package interfaces

type ModuleLogger interface {
	Info(content string)
	Warning(content string)
	Error(content string)
	Debug(content string)
}
