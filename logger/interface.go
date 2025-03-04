package logger

import "sync"

type LoggerEpic interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Trace(msg string, args ...interface{})
}

var (
	Logger LoggerEpic
	once   sync.Once
)

// func init() {
// 	once.Do(func() {
// 		Logger = NewLogrus() // Use logrus as default log, use default file path
// 	})
// }

func SetLogger(customLogger LoggerEpic) {
	once.Do(func() {
		Logger = customLogger
	})
}
