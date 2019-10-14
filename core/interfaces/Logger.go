package interfaces

import "io"

type Logger interface {
	Info(category string, format string, a ...interface{})
	Warn(category string, format string, a ...interface{})
	Error(category string, format string, a ...interface{})
	Debug(category string, format string, a ...interface{})
	io.Writer
}
