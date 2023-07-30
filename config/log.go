package config

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(c Conf) *zap.SugaredLogger {
	LogMode := zapcore.InfoLevel
	writeSyncer := getWriteSyncer(c)

	if c.App.Debug {
		writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout))
		LogMode = zapcore.DebugLevel
	}

	core := zapcore.NewCore(getEncoder(c), writeSyncer, LogMode)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	field := zap.Fields(zap.String("app_name", c.App.Name))
	sugarLogger := zap.New(core, caller, development, field).Sugar()
	// 封装自己的日志打印方法
	logger := sugarLogger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar()

	return logger
}

func getEncoder(c Conf) zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.TimeKey = "time"
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Local().Format(time.DateTime))
	}

	return zapcore.NewJSONEncoder(encodeConfig)
}

func getWriteSyncer(c Conf) zapcore.WriteSyncer {
	lumberjackSync := &lumberjack.Logger{
		Filename:   c.Zap.Dir + time.Now().Format(time.DateOnly) + ".log",
		MaxSize:    c.Zap.MaxSize,
		MaxBackups: c.Zap.MaxBackups,
		Compress:   c.Zap.Compress,
	}

	return zapcore.AddSync(lumberjackSync)
}
