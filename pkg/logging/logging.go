package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var DefaultLogger = InitializeDefaultLogger()

func InitializeDefaultLogger() *logrus.Logger {
	logger := logrus.New()
	envvar := os.Getenv("HAWKV6_DEBUG")
	if envvar != "" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return logger
}
