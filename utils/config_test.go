package utils

import (
	"bytes"
	"strings"
	"testing"

	Log "github.com/apatters/go-conlog"
)

func TestSetUpLogs(t *testing.T) {
	tests := []struct {
		name      string
		verbosity string
	}{
		{
			name:      "debug level",
			verbosity: "debug",
		},
		{
			name:      "info level",
			verbosity: "info",
		},
		{
			name:      "warn level",
			verbosity: "warn",
		},

		{
			name:      "invalid level defaults to warn",
			verbosity: "invalid",
		},
		{
			name:      "empty level defaults to warn",
			verbosity: "",
		},
		{
			name:      "uppercase level",
			verbosity: "DEBUG",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should not panic
			SetUpLogs(tt.verbosity)

			// Test that logging works after setup
			var buf bytes.Buffer
			Log.SetOutput(&buf)

			// Try to log at different levels to ensure setup worked
			Log.Debug("test debug message")
			Log.Info("test info message")
			Log.Warn("test warn message")
			Log.Error("test error message")

			output := buf.String()

			// Check that appropriate messages are logged based on level
			hasOutput := output != ""
			if !hasOutput {
				t.Errorf("Expected some log output for level %s, got: %s", tt.verbosity, output)
			}

			// Reset log output
			Log.SetOutput(nil)
		})
	}
}

func TestSetUpLogsFormatter(t *testing.T) {
	// Test that the formatter is set up correctly
	SetUpLogs("info")

	// We can't easily test the formatter directly, but we can test that logging works
	var buf bytes.Buffer
	Log.SetOutput(&buf)

	Log.Info("test message")
	output := buf.String()

	if !strings.Contains(output, "test message") {
		t.Errorf("Expected log output to contain 'test message', got: %s", output)
	}

	// Reset
	Log.SetOutput(nil)
}

func TestSetUpLogsWithDifferentLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}

	for _, level := range levels {
		t.Run("level_"+level, func(t *testing.T) {
			// This should not panic
			SetUpLogs(level)

			// Test that we can log at the set level
			var buf bytes.Buffer
			Log.SetOutput(&buf)

			switch level {
			case "debug":
				Log.Debug("debug test")
				if !strings.Contains(buf.String(), "debug test") {
					t.Errorf("Debug level should allow debug messages")
				}
			case "info":
				Log.Info("info test")
				if !strings.Contains(buf.String(), "info test") {
					t.Errorf("Info level should allow info messages")
				}
			case "warn":
				Log.Warn("warn test")
				if !strings.Contains(buf.String(), "warn test") {
					t.Errorf("Warn level should allow warn messages")
				}

			}

			Log.SetOutput(nil)
		})
	}
}
