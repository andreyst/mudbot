package botutil

import (
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LogConfig(logLevel zap.AtomicLevel, coloredLevels bool) zap.Config {
	logConfig := zap.NewProductionConfig()
	logConfig.DisableCaller = true
	logConfig.DisableStacktrace = true
	logConfig.Encoding = "console"
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	if !coloredLevels {
		logConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	logConfig.Level = logLevel
	logConfig.Sampling = nil

	return logConfig
}

func NewLogger(component string) *zap.SugaredLogger {
	logLevel, hasCommonLogLevel := os.LookupEnv("LOG_LEVEL")
	if !hasCommonLogLevel {
		logLevel = "INFO"
	}

	componentLogLevel, hasComponentLogLevel := os.LookupEnv("LOG_LEVEL_" + strings.ToUpper(component))
	if hasComponentLogLevel {
		logLevel = componentLogLevel
	}

	coloredLevels := true
	if strings.HasSuffix(os.Args[0], ".test") {
		logLevel = "DEBUG"
		coloredLevels = false
	}

	var zapLogLevel zap.AtomicLevel
	switch logLevel {
	case "DEBUG":
		zapLogLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "WARN":
		zapLogLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "INFO":
		zapLogLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "ERROR":
		zapLogLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		log.Fatalf("Unknown log level: %s", logLevel)
	}

	desugaredLogger, err := LogConfig(zapLogLevel, coloredLevels).Build()
	if err != nil {
		log.Fatalf("error while creating logger: %v\n", err)
	}

	return desugaredLogger.Sugar().Named(component)
}
