package main

import (
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	{
		var L *zap.Logger
		var S *zap.SugaredLogger

		zapConfig := zap.NewProductionConfig()
		if os.Getenv("mode") == "Dev" {
			zapConfig = zap.NewDevelopmentConfig()
		}
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		zapConfig.OutputPaths = []string{"stdout", "./" + time.Now().Format("2006-01-02") + ".log"}
		L, _ = zapConfig.Build()
		S = L.Sugar()

		L.Info("ordinary testing...", zap.String("type", "zap-logger"))
		S.Infow("ordinary testing...", "type", "zap-sugarlogger")
	}

	{
		// zapcore
		syncer := zapcore.AddSync(io.MultiWriter(
			os.Stdout,
			&lumberjack.Logger{
				Filename: "./" + time.Now().Format("2006-01-02") + ".log",
				MaxSize:  1, // 1MB
				MaxAge:   1, // 1Day
				// Compress: true,
			}))

		// encoder
		encoderConfig := zap.NewProductionEncoderConfig()
		if os.Getenv("mode") == "Dev" {
			encoderConfig = zap.NewDevelopmentEncoderConfig()
		}
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		encoder := zapcore.NewJSONEncoder(encoderConfig)

		//core
		core := zapcore.NewCore(encoder, syncer, zapcore.DebugLevel)

		S := zap.New(core, zap.WithCaller(true), zap.AddCallerSkip(0)).
			Sugar().Named("zap-demo").With("demo", "zap")

		for {
			S.Infow("customer testing...", "type", "zap-sugarlogger")
		}
	}
}
