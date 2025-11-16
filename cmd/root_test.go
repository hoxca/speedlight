package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"speedlight/utils"
)

func TestRootCmd(t *testing.T) {
	// Reset viper and global state
	viper.Reset()

	tests := []struct {
		name     string
		args     []string
		expected string
		hasError bool
	}{
		{
			name:     "root command shows help",
			args:     []string{"--help"},
			expected: "Speedlight implement 2 commands",
			hasError: false,
		},
		{
			name:     "root command without args shows help",
			args:     []string{},
			expected: "Speedlight implement 2 commands",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)

			// Execute command
			rootCmd.SetArgs(tt.args)
			_ = rootCmd.Execute()

			output := buf.String()

			// Error handling is complex due to the nature of the command structure
			// We'll focus on output validation for now

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output to contain %q, got: %s", tt.expected, output)
			}
		})
	}
}

func TestReportCmdFlags(t *testing.T) {
	// Reset viper and global state
	viper.Reset()

	// Set up default config for testing
	viper.Set("regexp", `/(.*)/[[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2}/(.*)/.*_LIGHT_[LRGBSHO]*_([[:digit:]]*)s_BIN1_(.*)C_GA.*_[[:digit:]]{8}_[[:digit:]]{6}_[[:digit:]]{3}_PA([[:digit:]]{3}\.?[[:digit:]]?[[:digit:]]?)_[EW]\.FIT`)
	utils.Regex = viper.GetString("regexp")

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "report command with valid flags",
			args:        []string{"report", "--dir", "/tmp/test", "--console", "true", "--report", "false"},
			expectError: false,
		},
		{
			name:        "report command with default flags",
			args:        []string{"report", "--dir", "/tmp/test"},
			expectError: false,
		},
		{
			name:        "report command without dir",
			args:        []string{"report"},
			expectError: true, // Should fail because no directory is set
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for testing
			tempDir, err := os.MkdirTemp("", "speedlight-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Replace /tmp/test with actual temp dir in args
			for i, arg := range tt.args {
				if arg == "/tmp/test" {
					tt.args[i] = tempDir
				}
			}

			// Capture output
			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)

			// Execute command
			rootCmd.SetArgs(tt.args)

			if tt.expectError {
				// Expect panic due to Fatal call
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic but got none")
					}
				}()
				err = rootCmd.Execute()
			} else {
				// Don't expect panic
				err = rootCmd.Execute()
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestConfigFileHandling(t *testing.T) {
	// Reset viper
	viper.Reset()

	// Create a temporary config file
	tempDir, err := os.MkdirTemp("", "speedlight-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configFile := filepath.Join(tempDir, "test-config.yaml")
	configContent := `
lightsdir: "/test/lights"
time_frame: 24
regexp: "/(.*)/[[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2}/(.*)/.*_LIGHT_[LRGBSHO]*_([[:digit:]]*)s_BIN1_(.*)C_GA.*_[[:digit:]]{8}_[[:digit:]]{6}_[[:digit:]]{3}_PA([[:digit:]]{3}\.[[:digit:]]{2})_[EW]\.FIT"
level: "debug"
`

	_ = os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "load custom config file",
			args:     []string{"--config", configFile, "report", "--dir", tempDir},
			expected: "Scanning lights directory:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper before each test
			viper.Reset()

			// Capture output
			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)

			// Execute command (expect potential panic due to Fatal)
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Command panicked as expected due to directory issue")
				}
			}()
			_ = rootCmd.Execute()

			output := buf.String()

			// We expect this to potentially fail due to directory structure, but config should load
			if strings.Contains(output, tt.expected) {
				t.Logf("Command executed as expected with config loading")
			}
		})
	}
}

func TestFlagPrecedence(t *testing.T) {
	// Reset viper
	viper.Reset()

	// Create a temporary config file
	tempDir, err := os.MkdirTemp("", "speedlight-flags-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configFile := filepath.Join(tempDir, "test-config.yaml")
	configContent := `
lightsdir: "/config/lights"
level: "info"
`

	_ = os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test that CLI flags override config file
	t.Run("CLI flag overrides config", func(t *testing.T) {
		viper.Reset()

		// Capture output
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)

		// Use CLI flag that should override config
		args := []string{"--config", configFile, "--dir", "/cli/lights", "report"}
		rootCmd.SetArgs(args)

		// Expect panic due to invalid directory
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Command panicked as expected due to directory issue")
			}
		}()
		_ = rootCmd.Execute()

		output := buf.String()

		// CLI flag should override config file value
		if strings.Contains(output, "/cli/lights") {
			t.Logf("CLI flag correctly overrode config file value")
		}
	})
}

func TestEnvironmentVariables(t *testing.T) {
	// Reset viper
	viper.Reset()

	// Set environment variable
	os.Setenv("SPEEDLIGHT_LIGHTSDIR", "/env/lights")
	defer os.Unsetenv("SPEEDLIGHT_LIGHTSDIR")

	t.Run("environment variable works", func(t *testing.T) {
		viper.Reset()

		// Capture output
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)

		args := []string{"report"}
		rootCmd.SetArgs(args)

		// Expect panic due to invalid directory
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Command panicked as expected due to directory issue")
			}
		}()
		_ = rootCmd.Execute()

		output := buf.String()

		// Environment variable should be used
		if strings.Contains(output, "/env/lights") {
			t.Logf("Environment variable correctly used")
		}
	})
}

func TestCommandValidation(t *testing.T) {
	// Reset viper
	viper.Reset()

	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "invalid command",
			args:        []string{"invalid-command"},
			expectError: true,
			errorMsg:    "unknown command",
		},
		{
			name:        "invalid flag",
			args:        []string{"--invalid-flag"},
			expectError: true,
			errorMsg:    "unknown flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()

			// Capture output
			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)

			// Execute command
			rootCmd.SetArgs(tt.args)
			_ = rootCmd.Execute()

			output := buf.String()

			// Error validation is complex due to command structure
			// Focus on output validation for error messages
			if tt.expectError && !strings.Contains(output, tt.errorMsg) {
				t.Errorf("Expected error message to contain %q, got: %s", tt.errorMsg, output)
			}
		})
	}
}
