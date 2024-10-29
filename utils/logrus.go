package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	logrus.TextFormatter
}

var (
	Log     *logrus.Logger
	LogJson *logrus.Logger
)

func init() {
	Log = logrus.New()
	LogJson = logrus.New()

	// Set logger to use the custom text formatter
	Log.SetFormatter(
		&logrus.TextFormatter{
			TimestampFormat: "2006/01/02 15:04:00",
			FullTimestamp:   true,
			ForceColors:     true,
			PadLevelText:    true,
		},
	)

	LogJson.SetFormatter(
		&logrus.JSONFormatter{
			TimestampFormat: "2006/01/02 15:04:00",
			PrettyPrint:     true,
		},
	)

	Log.SetOutput(os.Stdout)
}
