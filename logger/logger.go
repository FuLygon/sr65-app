package logger

import "github.com/sirupsen/logrus"

var Log = logrus.New()

func LogError(level logrus.Level, message string, err error) {
	entry := Log.WithField("error", err)
	switch level {
	case logrus.PanicLevel:
		entry.Panic(message)
	case logrus.FatalLevel:
		entry.Fatal(message)
	case logrus.ErrorLevel:
		entry.Error(message)
	case logrus.WarnLevel:
		entry.Warn(message)
	case logrus.InfoLevel:
		entry.Info(message)
	case logrus.DebugLevel:
		entry.Debug(message)
	case logrus.TraceLevel:
		entry.Trace(message)
	}
}

func LogErrorEmbed(message string, err error) {
	Log.WithField("error", err).Error(message)
	Log.Warn("error occuring when extracting embedded binaries, fallback to system binaries")
}
