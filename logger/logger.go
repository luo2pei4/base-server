package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *logrus.Logger

func init() {

	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(&lumberjack.Logger{
		Filename: "./base-server.log",
		MaxSize:  10,
		MaxAge:   10,
		Compress: true,
	})
}

func Infof(format string, args ...string) {
	logger.Infof(format, args)
}

func Info(args ...string) {
	logger.Info(args)
}

func Panicf(format string, args ...string) {
	logger.Panicf(format, args)
}
