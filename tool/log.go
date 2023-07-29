package tool

import (
	"errors"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/ozeer/sloth/config"
	"github.com/sirupsen/logrus"
)

var isTerm bool

func init() {
	isTerm = isatty.IsTerminal(os.Stdout.Fd())
}

var (
	LogAccess *logrus.Logger
	LogError  *logrus.Logger
)

func InitLog(c config.Conf) error {
	var err error

	// init logger
	LogAccess = logrus.New()
	LogError = logrus.New()

	if !isTerm {
		LogAccess.SetFormatter(&logrus.JSONFormatter{})
		LogError.SetFormatter(&logrus.JSONFormatter{})
	} else {
		LogAccess.Formatter = &logrus.TextFormatter{
			TimestampFormat: "2006/01/02 15:04:05",
			FullTimestamp:   true,
		}

		LogError.Formatter = &logrus.TextFormatter{
			TimestampFormat: "2006/01/02 15:04:05",
			FullTimestamp:   true,
		}
	}

	// set logger
	if err = SetLogLevel(LogAccess, c.Log.AccessLevel); err != nil {
		return errors.New("Set access log level error: " + err.Error())
	}

	if err = SetLogLevel(LogError, c.Log.ErrorLevel); err != nil {
		return errors.New("Set error log level error: " + err.Error())
	}

	if err = SetLogOut(LogAccess, c.Log.AccessLog); err != nil {
		return errors.New("Set access log path error: " + err.Error())
	}

	if err = SetLogOut(LogError, c.Log.ErrorLog); err != nil {
		return errors.New("Set error log path error: " + err.Error())
	}

	return nil
}

func SetLogOut(log *logrus.Logger, outString string) error {
	switch outString {
	case "stdout":
		log.Out = os.Stdout
	case "stderr":
		log.Out = os.Stderr
	default:
		f, err := os.OpenFile(outString, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)

		if err != nil {
			return err
		}

		log.Out = f
	}

	return nil
}

func SetLogLevel(log *logrus.Logger, levelString string) error {
	level, err := logrus.ParseLevel(levelString)

	if err != nil {
		return err
	}

	log.Level = level

	return nil
}
