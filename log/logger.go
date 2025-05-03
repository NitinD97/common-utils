package log

import (
	"fmt"
	"github.com/OneOrbit/common-utils/configuration"
	"github.com/OneOrbit/common-utils/context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zap.Logger
}

var logger *Logger

func NewLogger(ctx *context.Context) *zap.Logger {
	config := configuration.GetConfig()
	environment := config.GetString("environment")
	if environment != "development" {
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.GetString("log.file_name"),
			MaxSize:    config.GetInt("log.max_size"), // megabytes
			MaxBackups: config.GetInt("log.max_backups"),
			MaxAge:     config.GetInt("log.max_age"), // days
			Compress:   true,
		})
		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncodeDuration = zapcore.MillisDurationEncoder
		core := zapcore.NewTee(zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg),
			writer,
			zap.DebugLevel,
		))
		logger = &Logger{
			zap.New(core),
		}
	} else {
		l, err := zap.NewDevelopment()

		if err != nil {
			panic(fmt.Errorf("unable to initialize utils\n %w", err))
		}
		logger = &Logger{
			Logger: l,
		}

	}
	logger.Logger = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(0), zap.AddStacktrace(zap.ErrorLevel))
	logger.Logger = logger.WithContext(ctx)

	return logger.Logger
}

func GetLogger() *zap.Logger {
	return logger.Logger
}

func (logger *Logger) WithContext(ctx *context.Context) *zap.Logger {
	fields := make([]zap.Field, 0)
	for k, v := range ctx.GetContextMap() {
		fields = append(fields, zap.String(k, v))
	}

	return logger.Logger.With(fields...)
}
