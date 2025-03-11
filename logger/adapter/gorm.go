package adapter

import (
	"context"
	"time"

	coreLogger "github.com/epicconsult/pkgep/logger"

	"gorm.io/gorm/logger"
)

type GormResponse struct {
	Sql  string `json:"sql"`
	Rows int64  `json:"rows"`
	Time string `json:"time"`
}

type GormLogger struct {
	epicLogger coreLogger.LoggerEpic
	logLevel   logger.LogLevel
}

func NewGormLogger(epicLogger coreLogger.LoggerEpic) logger.Interface {
	return &GormLogger{
		epicLogger: epicLogger,
		logLevel:   logger.Info,
	}
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	g.logLevel = level
	return g
}

func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	g.epicLogger.Info(ctx, msg, data...)
}

func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	g.epicLogger.Warn(ctx, msg, data...)
}

func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	g.epicLogger.Error(ctx, msg, data...)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// Capture SQL, number of rows
	sql, rows := fc()
	duration := time.Since(begin)

	// fields := map[string]interface{}{
	// 	"sql":  sql,
	// 	"rows": rows,
	// 	"time": duration.String(),
	// }

	resp := GormResponse{
		Sql:  sql,
		Rows: rows,
		Time: duration.String(),
	}

	logCtx := context.WithValue(context.Background(), coreLogger.LogHeader, resp)

	// Log response.
	if err != nil {
		g.epicLogger.InfoWithAction(logCtx, coreLogger.DBRESPONSE, "GORM Query Error: "+err.Error())
	} else {
		g.epicLogger.InfoWithAction(logCtx, coreLogger.DBRESPONSE, "GORM Query Executed Successfully")
	}
}
