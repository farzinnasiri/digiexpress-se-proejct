package log

import (
	"github.com/sirupsen/logrus"
)

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Info(args ...interface{}) {
	logrus.Infoln(args...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Debug(args ...interface{}) {
	logrus.Debugln(args...)
}

func Error(msg string, err error) {
	logrus.WithError(err).Errorln(msg)
}

func Panic(msg string, err error) {
	logrus.WithError(err).Panicln(msg)
}

func Warning(msg string) {
	logrus.Warning(msg)
}
