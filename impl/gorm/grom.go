package gorm

import (
	"context"
	"errors"
	"time"

	"github.com/thoohv5/logger"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	ds "github.com/thoohv5/orm/standard"
)

// gDb
type gDb struct {
	*ds.Options
	db *gorm.DB
}

func NewGDb() ds.IConnect {
	return &gDb{}
}

func CopyGDb(gdb *gorm.DB, sos ...ds.ServerOption) ds.IBuilder {
	opts := new(ds.Options)
	for _, so := range sos {
		so(opts)
	}

	if opts.GetIsSetLog() {
		gdb.Logger = WithLogger(opts.GetLogger(), opts.GetLoggerLevel())
	}

	return &gDb{
		Options: opts,
		db:      gdb,
	}
}

func (g *gDb) gDB() *gorm.DB {
	return g.db
}

func (g *gDb) Write() ds.IBuilder {
	// TODO implement me
	panic("implement me")
}

func (g *gDb) Read() ds.IBuilder {
	// TODO implement me
	panic("implement me")
}

type defaultLogger struct {
	logger.ILogger
	Level                     gormlogger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

func WithLogger(logger logger.ILogger, level string) gormlogger.Interface {
	return &defaultLogger{
		ILogger:                   logger,
		Level:                     parseLevel(level),
		SlowThreshold:             100 * time.Millisecond,
		IgnoreRecordNotFoundError: false,
	}
}

func (d *defaultLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	d.Level = level
	return d
}

func (d *defaultLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if d.Level <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && d.Level >= gormlogger.Error && (!d.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		d.ILogger.Error("trace", map[string]interface{}{
			// "ctx": ctx,
			"err":     err.Error(),
			"elapsed": elapsed,
			"rows":    rows,
			"sql":     sql,
		})
	case d.SlowThreshold != 0 && elapsed > d.SlowThreshold && d.Level >= gormlogger.Warn:
		sql, rows := fc()
		d.ILogger.Warn("trace", map[string]interface{}{
			// "ctx":     ctx,
			"err":     err.Error(),
			"elapsed": elapsed,
			"rows":    rows,
			"sql":     sql,
		})
	case d.Level >= gormlogger.Info:
		sql, rows := fc()
		d.ILogger.Info("trace", map[string]interface{}{
			// "ctx":     ctx,
			"err":     err.Error(),
			"elapsed": elapsed,
			"rows":    rows,
			"sql":     sql,
		})
	}
}

func (d *defaultLogger) Info(ctx context.Context, s string, i ...interface{}) {
	// TODO implement me
	panic("implement me")
}

func (d *defaultLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	// TODO implement me
	panic("implement me")
}

func (d *defaultLogger) Error(ctx context.Context, s string, i ...interface{}) {
	// TODO implement me
	panic("implement me")
}

// 日志类别: debug, warn, info，error
func parseLevel(level string) gormlogger.LogLevel {
	zl := gormlogger.Silent
	switch level {
	case "debug":
		zl = gormlogger.Warn
	case "warn":
		zl = gormlogger.Warn
	case "info":
		zl = gormlogger.Info
	case "error":
		zl = gormlogger.Error
	}
	return zl
}
