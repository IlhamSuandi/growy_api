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

	Log.SetFormatter(
		&logrus.TextFormatter{
			TimestampFormat: "2006/01/02 15:04",
			FullTimestamp:   true,
			ForceColors:     true,
			PadLevelText:    true,
		},
	)

	Log.SetOutput(os.Stdout)
}
