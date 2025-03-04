package logger

import (
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

type LogrusEpic struct {
	logger     *logrus.Logger
	rotation   RotationType
	maxSize    int // work only with rotationType set to Timestamp
	maxBackups int // work only with rotationType set to Timestamp
	appName    string
	path       string // default root running
}

// ** Support configuration via "functional options pattern"
type LogrusOption func(*LogrusEpic)

func WithMaxSize(maxSize int) LogrusOption {
	return func(le *LogrusEpic) {
		le.maxSize = maxSize
	}
}

func WithMaxBackups(maxBackups int) LogrusOption {
	return func(le *LogrusEpic) {
		le.maxBackups = maxBackups
	}
}

func WithRotationType(rotationType RotationType) LogrusOption {
	return func(le *LogrusEpic) {
		le.rotation = rotationType
	}
}

func WithAppName(name string) LogrusOption {
	return func(le *LogrusEpic) {
		le.appName = name
	}
}

func WithPath(path string) LogrusOption {
	return func(le *LogrusEpic) {
		le.path = path
	}
}

func NewLogrus(options ...LogrusOption) *LogrusEpic {

	// Default configuration.
	client := LogrusEpic{
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

func (l *LogrusEpic) Info(msg string, args ...interface{}) {
	l.logger.Info(append([]interface{}{msg}, args...)...)
}

func (l *LogrusEpic) Error(msg string, args ...interface{}) {
	l.logger.Error(append([]interface{}{msg}, args...)...)
}

func (l *LogrusEpic) Warn(msg string, args ...interface{}) {
	l.logger.Warn(append([]interface{}{msg}, args...)...)
}

func (l *LogrusEpic) Trace(msg string, args ...interface{}) {
	l.logger.Trace(append([]interface{}{msg}, args...)...)
}
