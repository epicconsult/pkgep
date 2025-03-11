package pkgep

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type Logger struct {
	*logrus.Logger
	TransactionID string
	DeviceID      string
	UserID        string
	ParID         string
	Role          string
	AuthType      string
	Path          string
	Method        string
}
type LogAction int

const (
	INBOUND      LogAction = iota
	OUTBOUND     LogAction = iota
	DBRESPOSNE   LogAction = iota
	HTTPREQEST   LogAction = iota
	HTTPRESPONSE LogAction = iota
)

var logr *Logger

func SetupLogger(appname string) {

	logr = &Logger{
		Logger:        logrus.New(),
		TransactionID: uuid.New().String(),
	}

	// Log as JSON instead of the default ASCII formatter.
	logr.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})
	logr.WithFields(logrus.Fields{
		"app": appname,
	})

	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0770)
		// panic(err)
	}

	date := time.Now().UTC()
	logFile, err := os.OpenFile("logs/log_"+date.Format("01-02-2006_15")+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	logr.SetOutput(mw)

}

func MockLogger(appname string) *Logger {
	logr = &Logger{
		Logger: logrus.New(),
	}
	logr.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})
	logr.WithFields(logrus.Fields{
		"app": appname,
	})
	return logr
}

func NewLogger() *Logger {
	return logr
}

func SetHeaderLog(token *VerifiedToken, c *fiber.Ctx) {
	logr.DeviceID = token.Sub.DeviceID
	logr.UserID = strconv.Itoa(token.Sub.UserID)
	logr.ParID = strconv.Itoa(token.Sub.ParID)
	logr.Role = token.Sub.Role
	logr.AuthType = token.Sub.AuthType
	logr.Path = c.Path()
	logr.Method = c.Method()
}



func (l *Logger) LogInformation(Action LogAction, args ...interface{}) {
	componentName := strings.Split(l.Path, "/")[len(strings.Split(l.Path, "/"))-1]
	if Action == INBOUND {
		l.Logger.WithFields(logrus.Fields{
			"action":            "[INBOUND]",
			"actionDescription": fmt.Sprintf("Start receiving request from API : %s", l.Path),
			"componentName":     componentName,
			"methodName":        l.Method,
			"deviceID":          l.DeviceID,
			"userID":            l.UserID,
			"parID":             l.ParID,
			"deviceName":        l.AuthType,
			"transactionID":     l.TransactionID,
		}).Info(args...) // Add ... after args
	} else if Action == OUTBOUND {
		l.Logger.WithFields(logrus.Fields{
			"action":            "[OUTBOUND]",
			"actionDescription": fmt.Sprintf("Responding to request from API : %s", l.Path),
			"componentName":     componentName,
			"methodName":        l.Method,
			"deviceID":          l.DeviceID,
			"userID":            l.UserID,
			"parID":             l.ParID,
			"deviceName":        l.AuthType,
			"transactionID":     l.TransactionID,
		}).Info(args...)
	} else if Action == HTTPREQEST {
		l.Logger.WithFields(logrus.Fields{
			"action":            "[HTTPREQEST]",
			"actionDescription": fmt.Sprintf("Start receiving request from API : %s", l.Path),
			"componentName":     componentName,
			"methodName":        l.Method,
			"deviceID":          l.DeviceID,
			"userID":            l.UserID,
			"parID":             l.ParID,
			"deviceName":        l.AuthType,
			"transactionID":     l.TransactionID,
		}).Info(args...)
	} else if Action == HTTPRESPONSE {
		l.Logger.WithFields(logrus.Fields{
			"action":            "[HTTPRESPONSE]",
			"actionDescription": fmt.Sprintf("Responding to request from API : %s", l.Path),
			"componentName":     componentName,
			"methodName":        l.Method,
			"deviceID":          l.DeviceID,
			"userID":            l.UserID,
			"parID":             l.ParID,
			"deviceName":        l.AuthType,
			"transactionID":     l.TransactionID,
		}).Info(args...)
	} else if Action == DBRESPOSNE {
		l.Logger.WithFields(logrus.Fields{
			"action":            "[DBRESPONSE]",
			"actionDescription": "End proccess.",
			"componentName":     componentName,
			"deviceID":          l.DeviceID,
			"userID":            l.UserID,
			"parID":             l.ParID,
			"deviceName":        l.AuthType,
			"transactionID":     l.TransactionID,
		}).Info(args...)
	}
}

func (l *Logger) LogError(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc()
	sqltype, stmt := "Command Log", ""
	componentName := strings.Split(l.Path, "/")[len(strings.Split(l.Path, "/"))-1]
	if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(sql)), "SELECT") {
		stmt = "SELECT"
		sqltype = "Query Log"
	} else if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(sql)), "INSERT") {
		stmt = "INSERT"
	} else if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(sql)), "UPDATE") {
		stmt = "UPDATE"
	} else if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(sql)), "DELETE") {
		stmt = "DELETE"
	}

	switch {
	case err != nil && l.Logger.IsLevelEnabled(logrus.ErrorLevel):
		l.Logger.WithFields(logrus.Fields{
			"action":            "[DBRESPONSE]",
			"actionDescription": fmt.Sprintf("Responding to request from API : %s", l.Path),
			"errorCode":         "500",
			"type":              strings.ToLower(stmt),
			"queryString":       sql,
			"transactionID":     l.TransactionID,
		}).Error(err)
	case l.Logger.IsLevelEnabled(logrus.InfoLevel):
		l.Logger.WithFields(logrus.Fields{
			"action":            "[DBREQUEST]",
			"actionDescription": fmt.Sprintf("Start %s proccessing.", stmt),
			"componentName":     componentName,
			"deviceID":          l.DeviceID,
			"userID":            l.UserID,
			"parID":             l.ParID,
			"deviceName":        l.AuthType,
			"transactionID":     l.TransactionID,
			"queryString":       sql,
		}).Info(sqltype)
	}

}

func (l Logger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := logrus.New()
	newLogger.SetLevel(logrus.Level(level))
	return &Logger{Logger: newLogger}
}

func (l Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.WithContext(ctx).Infof(msg, data...)
}

func (l Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.WithContext(ctx).Warnf(msg, data...)
}

func (l Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.WithContext(ctx).Errorf(msg, data...)
}
