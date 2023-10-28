package logger

import (
	"os"

	"github.com/lucyanddarlin/lucy-ez-admin/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type logger struct {
	*zap.Logger
	c *config.Log
}

type Logger interface {
	// WithID
	//
	//  @Description: 设置链路日志 id
	//  @param id: 请求唯一 id
	//  @return: *zap.Logger
	WithID(id string) *zap.Logger
	// Field
	//
	//  @Description: 设置链路日志 id
	//  @return: string
	Field() string
}

// New
//
//	@Description: 初始化日志器
//	@receiver: conf 日志配置
//	@receiver: srvName 服务名
//	@return: Logger 日志器
func New(conf *config.Log, srvName string) Logger {
	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "lever",
		NameKey:        "log",
		CallerKey:      "caller",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,                          // 小写编码器
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"), // ISO08601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 输出器配置
	var output []zapcore.WriteSyncer
	for _, val := range conf.Output {
		if val == "stdout" {
			output = append(output, zapcore.AddSync(os.Stdout))
		}
		if val == "file" {
			output = append(output, zapcore.AddSync(&lumberjack.Logger{
				Filename:   conf.File.Name,
				MaxSize:    conf.File.MaxSize,
				MaxBackups: conf.File.MaxBackup,
				MaxAge:     conf.File.MaxAge,
				Compress:   conf.File.Compress,
			}))
		}
	}

	// 创建 zap-core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(output...),
		zapcore.Level(conf.Level),
	)

	// 创建 zap
	return &logger{
		c:      conf,
		Logger: zap.New(core, zap.AddCaller(), zap.Fields(zap.String("service", srvName))),
	}
}

func (l *logger) Field() string {
	return l.c.Field
}

func (l *logger) WithID(id string) *zap.Logger {
	return l.Logger.With(zap.Any(l.c.Field, id))
}
