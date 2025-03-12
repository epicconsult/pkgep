package logger

import (
	"context"
	"sync"
)

type LogHeaderContextKey string

const (
	LogHeader LogHeaderContextKey = "header" //todo: provide helper function to update log header name.
)

type LogAction int

const (
	INBOUND LogAction = iota
	OUTBOUND
	DBREQUEST
	DBRESPONSE
)

// Array get allocated in stack, efficient and faster than heap alloc.
var logActionName = [...]string{
	"INBOUND",
	"OUTBOUND",
	"DBREQUEST",
	"DBRESPONSE",
}

// Core Logger Interface
type EpicLogger interface {
	Info(ctx context.Context, msg string, data ...any)
	Error(ctx context.Context, msg string, data ...any)
	Warn(ctx context.Context, msg string, data ...any)
	Trace(ctx context.Context, msg string, args ...any)

	InfoWithAction(ctx context.Context, action LogAction, msg string, data ...any)
}

var (
	Logger EpicLogger
	once   sync.Once
)

func SetLogger(customLogger EpicLogger) {
	once.Do(func() {
		Logger = customLogger
	})
}
