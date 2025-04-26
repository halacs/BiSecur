package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	log *logrus.Logger
)

func NewLogger() *logrus.Logger {
	if log != nil {
		return log
	}

	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		//ForceColors:   true,
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetReportCaller(false)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	return log
}
