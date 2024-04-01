package logger

import (
	"log/slog"
	"os"
)

var Logger = slog.Default()

type LogConfig struct {
	LogLevel   string `json:"logLevel"`
	OutputPath string `json:"outputPath"`
	OutputType string `json:"outputType"`
}

// TODO: log level
// TODO: log output
// TODO: log rotation
func InitLogger(LogConfig) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}
