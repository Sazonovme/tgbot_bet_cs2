package logger

import (
	"github.com/sirupsen/logrus"
)

// PanicLevel 0
// FatalLevel 1
// ErrorLevel 2
// WarnLevel 3
// InfoLevel 4
// DebugLevel 5
// TraceLevel 6

var logger *logrus.Logger

func InitLogger(level int, formatProduction bool) {
	logger = logrus.New()

	logger.SetLevel(logrus.Level(level))

	if formatProduction {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
	}
}

func Fatal(textMessage string, location string, err error) {
	logger.WithFields(logrus.Fields{
		"text":     textMessage,
		"location": location,
		"err":      err,
	}).Fatal()
}

func Error(textMessage string, location string, err error) {
	logger.WithFields(logrus.Fields{
		"text":     textMessage,
		"location": location,
		"err":      err,
	}).Error()
}

func Warn(textMessage string, location string, err error) {
	logger.WithFields(logrus.Fields{
		"text":     textMessage,
		"location": location,
		"err":      err,
	}).Warn()
}

func Info(textMessage string, location string, err error) {
	logger.WithFields(logrus.Fields{
		"text":     textMessage,
		"location": location,
		"err":      err,
	}).Info()
}

func Debug(textMessage string, location string, err error) {
	logger.WithFields(logrus.Fields{
		"text":     textMessage,
		"location": location,
		"err":      err,
	}).Debug()
}
