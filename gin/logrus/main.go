package main

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type AppHook struct {
	AppName string
}

func (h *AppHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *AppHook) Fire(entry *logrus.Entry) error {
	entry.Data["app"] = h.AppName
	return nil
}

func main() {
	logrus.AddHook(&AppHook{AppName: "awesome-web"})
	logrus.Info("info msg")

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	if os.Getenv("mode") == "Dev" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	logrus.SetOutput(io.MultiWriter(os.Stdout,
		&lumberjack.Logger{
			Filename: "logrus.log",
			MaxSize:  1,
			MaxAge:   1,
			Compress: true,
		}))

	for {
		logrus.WithField("ip", "127.0.0.1").Info("connecting...")

		logger := logrus.WithFields(
			logrus.Fields{
				"ip":   "localhost",
				"port": "8080",
			},
		)
		logger.Errorln("testing...")
	}
}
