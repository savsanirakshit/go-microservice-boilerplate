package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var ServiceLogger *logrus.Logger

var serviceLogFile *os.File

func ConfigLogger(logDir, logLevel string) {
	configureServiceLogger(logDir, logLevel)
}

func configureServiceLogger(logDir, logLevel string) {
	ServiceLogger = logrus.New()
	level, _ := logrus.ParseLevel(logLevel)
	ServiceLogger.SetLevel(level)
	serviceLogFile, _ = os.OpenFile(logDir+"/service.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	ServiceLogger.SetOutput(io.MultiWriter(os.Stdout, serviceLogFile))
}
