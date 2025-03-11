package logger

import (
	"context"
	"sync"
)

type LogHeaderContextKey string

const (
	LogHeader LogHeaderContextKey = "header" // provide helper function to update log header name.
)

type LogAction int

const (
	INBOUND LogAction = iota
	OUTBOUND
	DBREQUEST
	DBRESPONSE
)

var logActionName = [...]string{
	"INBOUND",
	"OUTBOUND",
	"DBREQUEST",
	"DBRESPONSE",
}

// Core Logger Interface
type LoggerEpic interface {
	Info(ctx context.Context, msg string, data ...any)
	Error(ctx context.Context, msg string, data ...any)
	Warn(ctx context.Context, msg string, data ...any)
	Trace(ctx context.Context, msg string, args ...any)

	InfoWithAction(ctx context.Context, action LogAction, msg string, data ...any)
}

var (
	Logger LoggerEpic
	once   sync.Once
)

func SetLogger(customLogger LoggerEpic) {
	once.Do(func() {
		Logger = customLogger
	})
}
