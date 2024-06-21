package util

import (
	"encoding/json"
	"log/slog"
)

func ToJson(data []byte) map[string]interface{} {
	// Unmarshal JSON into a generic map
	var responseMap map[string]interface{}
	json.Unmarshal(data, &responseMap)
	return responseMap
}

func LogLevelHandler(level string) *slog.HandlerOptions {
	var hOptions slog.HandlerOptions
	switch level {
	case "DEBUG":
		hOptions = slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	case "INFO":
		hOptions = slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	case "ERROR":
		hOptions = slog.HandlerOptions{
			Level: slog.LevelError,
		}
	default:
		hOptions = slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	}
	return &hOptions
}
