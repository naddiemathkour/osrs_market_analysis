package logging

import (
	_ "os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger = logrus.New()

func init() {
	// Set up log rotation with lumberjack
	Logger.SetOutput(&lumberjack.Logger{
		Filename:   "scheduler.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	})

	// Set log format to JSON for structured logging
	Logger.SetFormatter(&logrus.JSONFormatter{})
}
