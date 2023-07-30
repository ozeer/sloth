package global

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

// 错误日志打印方法
func Error(errMsg ...interface{}) {
	Logger.Error(errMsg)
	// os.Exit(1)
}

func Errorf(errMsg string, args ...interface{}) {
	Logger.Errorf(errMsg, args...)
	// os.Exit(1)
}

// 正常访问日志打印方法
func Info(infoMsg ...interface{}) {
	Logger.Info(infoMsg)
	// os.Exit(1)
}

// 正常访问日志格式化打印方法
func InfoF(infoMsg string, args ...interface{}) {
	Logger.Infof(infoMsg, args...)
	// os.Exit(1)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
	// os.Exit(1)
}
