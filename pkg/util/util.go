package util

import (
	"encoding/json"
	"log/slog"
)

// ToJson converts a JSON byte slice into a map[string]interface{}.
func ToJson(data []byte) map[string]interface{} {
	// Unmarshal JSON into a generic map
	var responseMap map[string]interface{}
	json.Unmarshal(data, &responseMap)
	return responseMap
}

// LogLevelHandler returns a slog.HandlerOptions pointer based on the provided log level string.
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
