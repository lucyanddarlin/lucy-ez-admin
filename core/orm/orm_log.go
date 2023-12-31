package orm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lucyanddarlin/lucy-ez-admin/config"
	coreLog "github.com/lucyanddarlin/lucy-ez-admin/core/log"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

// zap 适配 gorm 日志
type sqlLog struct {
	logger        coreLog.Logger
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

func newOrmLog(conf config.Orm, log coreLog.Logger) logger.Interface {
	return &sqlLog{
		logger:        log,
		LogLevel:      logger.LogLevel(conf.Level),
		SlowThreshold: conf.SlowThreshold,
	}
}

func (l *sqlLog) Log(ctx context.Context) *zap.Logger {
	traceId, _ := ctx.Value(l.logger.Field()).(string)
	return l.logger.WithID(traceId).WithOptions(zap.AddCallerSkip(3))
}

func (l *sqlLog) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *sqlLog) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Log(ctx).Info("SQL 信息", getSqlInfo("", fmt.Sprintf(msg, data...), 0, 0, false)...)
	}

}

// Warn implements logger.Interface.
func (l *sqlLog) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Log(ctx).Info("SQL 告警", getSqlInfo("", fmt.Sprintf(msg, data...), 0, 0, false)...)
	}
}

// Error implements logger.Interface.
func (l *sqlLog) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Log(ctx).Info("SQL 错误", getSqlInfo("", fmt.Sprintf(msg, data...), 0, 0, false)...)
	}
}

// Trace implements logger.Interface.
func (l *sqlLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	costTime := float64(elapsed.Abs().Nanoseconds()) / 1e6
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound)):
		sql, rows := fc()
		l.Log(ctx).Info("SQL 错误", getSqlInfo(err.Error(), sql, rows, costTime, false)...)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		l.Log(ctx).Info("SQL 告警", getSqlInfo("", sql, rows, costTime, true)...)
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		l.Log(ctx).Info("SQL 信息", getSqlInfo("", sql, rows, costTime, false)...)
	}

}

func getSqlInfo(err, sql string, rows int64, costTime float64, slow bool) []zap.Field {
	return []zap.Field{
		zap.String("err", err),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.String("time", fmt.Sprintf("%vms", costTime)),
		zap.Bool("slow", slow),
	}

}
