package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	Log.SetOutput(os.Stdout)

	Log.SetFormatter(&logrus.JSONFormatter{})

	Log.SetLevel(logrus.InfoLevel)
}

func Info(msg string, fields logrus.Fields) {
	Log.WithFields(fields).Info(msg)
}

func Warn(msg string, fields logrus.Fields) {
	Log.WithFields(fields).Warn(msg)
}

func Error(msg string, err error, fields logrus.Fields) {
	Log.WithError(err).
		WithFields(fields).
		Error(msg)
}
