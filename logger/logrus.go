package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type RotationType int

const (
	Date      RotationType = iota // default "yyyy-mm-dd"
	Timestamp                     // "yyyy-mm-ddT-hh-mm-ss"
)

type EpicLogrus struct {
	logger     *logrus.Logger
	rotation   RotationType
	maxSize    int // work only with rotationType set to Timestamp
	maxBackups int // work only with rotationType set to Timestamp
	appName    string
	path       string // default root running
	// infoFunc   func() // user can customize logging behavior
	// errFunc    func()
}

// ** Support configuration via "functional options pattern"
type LogrusOption func(*EpicLogrus)

func WithMaxSize(maxSize int) LogrusOption {
	return func(le *EpicLogrus) {
		le.maxSize = maxSize
	}
}

func WithMaxBackups(maxBackups int) LogrusOption {
	return func(le *EpicLogrus) {
		le.maxBackups = maxBackups
	}
}

func WithRotationType(rotationType RotationType) LogrusOption {
	return func(le *EpicLogrus) {
		le.rotation = rotationType
	}
}

func WithAppName(name string) LogrusOption {
	return func(le *EpicLogrus) {
		le.appName = name
	}
}

func WithPath(path string) LogrusOption {
	return func(le *EpicLogrus) {
		le.path = path
	}
}

func NewLogrus(options ...LogrusOption) *EpicLogrus {

	// Default configuration.
	client := EpicLogrus{
		rotation:   Date,
		maxSize:    500,
		maxBackups: 3,
		appName:    "epic-app",
	}

	// Apply custom logrus configuration
	for _, option := range options {
		option(&client)
	}

	client.logger = logrus.New()

	if client.rotation == Timestamp {
		client.logger.SetOutput(&lumberjack.Logger{
			Filename:   fmt.Sprintf("logs/%s.log", client.appName),
			MaxSize:    client.maxSize, // megabytes
			MaxBackups: client.maxBackups,
			MaxAge:     28, //days
			Compress:   true,
		})
	} else {
		// Create parent logs dir
		if _, err := os.Stat("logs"); os.IsNotExist(err) {
			err := os.Mkdir("logs", 0770)
			panic(err)
		}

		date := time.Now().UTC()
		logFile, err := os.OpenFile("logs/log_"+date.Format("01-02-2006_15")+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		client.logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}

	// Log detail configuration.
	client.logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})
	client.logger.WithFields(logrus.Fields{
		"app": client.appName,
	})

	return &client
}

func (l *EpicLogrus) Info(ctx context.Context, msg string, data ...any) {
	// Process context to be log as metadata
	fields := structToJson(ctx.Value(LogHeader))

	l.logger.WithFields(fields).Info(append([]any{msg}, data...)...)
}

func (l *EpicLogrus) Error(ctx context.Context, msg string, data ...any) {
	fields := structToJson(ctx.Value(LogHeader))
	l.logger.WithFields(fields).Error(append([]any{msg}, data...)...)
}

func (l *EpicLogrus) Warn(ctx context.Context, msg string, data ...any) {
	fields := structToJson(ctx.Value(LogHeader))
	l.logger.WithFields(fields).Warn(append([]any{msg}, data...)...)
}

func (l *EpicLogrus) Trace(ctx context.Context, msg string, data ...any) {
	fields := structToJson(ctx.Value(LogHeader))
	l.logger.WithFields(fields).Trace(append([]any{msg}, data...)...)
}

func (l *EpicLogrus) InfoWithAction(ctx context.Context, action LogAction, msg string, data ...any) {
	actionName := logActionName[action]
	fields := structToJson(ctx.Value(LogHeader))
	fields["action"] = actionName
	l.logger.WithFields(fields).Info(append([]any{msg}, data...)...)
}

// convert any struct to map.
func structToJson(st any) map[string]any {
	data, _ := json.Marshal(st)
	var kv map[string]any
	json.Unmarshal(data, &kv)
	return kv
}
