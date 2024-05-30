package logger

import "github.com/sirupsen/logrus"

var log = logrus.New()

func Info(message string) {
	log.Info(message)
}

func Warn(message string, err ...error) {
	if len(err) > 0 {
		log.WithError(err[0]).Warn(message)
	} else {
		log.Warn(message)
	}
}

func Error(message string, err ...error) {
	if len(err) > 0 {
		log.WithError(err[0]).Error(message)
	} else {
		log.Error(message)
	}
}

func Fatal(message string, err ...error) {
	if len(err) > 0 {
		log.WithError(err[0]).Panic(message)
	} else {
		log.Fatal(message)
	}
}
