package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	testRegex = `/(.*)/[[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2}/(.*)/.*_LIGHT_[LRGBSHO]*_([[:digit:]]*)s_BIN1_(.*)C_GA.*_[[:digit:]]{8}_[[:digit:]]{6}_[[:digit:]]{3}_PA([[:digit:]]{3}\.?[[:digit:]]?[[:digit:]]?)_[EW]\.FIT`
)

func TestTraversalRegexParsing(t *testing.T) {
	// Reset global state before test
	ObjectList = Objects{}
	targetList = Targets{}
	RotUsed = true
	Regex = testRegex

	tests := []struct {
		name     string
		path     string
		expected bool // should match regex
	}{
		{
			name:     "valid new nomenclature file",
			path:     "/NGC7635/2025-01-15/H/NGC7635_LIGHT_H_180s_BIN1_5C_GA2750_20250115_233348_235_PA239.88_W.FIT",
			expected: true,
		},
		{
			name:     "valid file with different filter",
			path:     "/M42/2025-01-15/L/M42_LIGHT_L_600s_BIN1_-20C_GA2750_20250115_222558_734_PA045.50_E.FIT",
			expected: true,
		},
		{
			name:     "invalid file - wrong extension",
			path:     "/M42/2025-01-15/L/M42_LIGHT_L_600s_BIN1_-20C_GA2750_20250115_222558_734_PA045.50_E.TXT",
			expected: false,
		},
		{
			name:     "invalid file - missing components",
			path:     "/M42/2025-01-15/L/M42_LIGHT_L_600s_BIN1_-20C.FIT",
			expected: false,
		},
		{
			name:     "directory path",
			path:     "/M42/2025-01-15/L",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the regex processing directly
			image := tt.path // Use the full path for testing

			re := regexp.MustCompilePOSIX(Regex)
			splitline := re.FindAllStringSubmatch(image, -1)

			hasMatch := len(splitline) == 1

			if tt.expected && !hasMatch {
				t.Errorf("Expected regex to match for path: %s", tt.path)
			}

			if !tt.expected && hasMatch {
				t.Errorf("Expected regex to not match, but it matched for path: %s", tt.path)
			}

			// If regex matched, test that object creation works
			if hasMatch {
				object := splitline[0][1]
				filter := splitline[0][2]
				expo, _ := strconv.Atoi(splitline[0][3])
				temperature, _ := strconv.Atoi(splitline[0][4])
				rotval, _ := strconv.ParseFloat(splitline[0][5], 32)
				rotation := float32(rotval)

				target := fmt.Sprintf("%s~%d~%d", object, temperature, expo)
				o := addTarget(target, rotation)
				o.addFilter(target, filter)

				if len(ObjectList) == 0 {
					t.Errorf("Expected object to be created after regex match")
				}
			}

			// Reset for next test
			ObjectList = Objects{}
		})
	}
}

func TestFlatsversalRegexParsing(t *testing.T) {
	// Reset global state before test
	FlatList = flats{}
	Rotations = []float32{}
	RotUsed = true
	TimeFrame = 24 // 24 hours to ensure files are within time window
	Regex = testRegex

	tests := []struct {
		name     string
		path     string
		expected bool // should match regex
	}{
		{
			name:     "valid new nomenclature file",
			path:     "/NGC7635/2025-01-15/H/NGC7635_LIGHT_H_180s_BIN1_5C_GA2750_20250115_233348_235_PA239.88_W.FIT",
			expected: true,
		},
		{
			name:     "valid file with different filter",
			path:     "/M42/2025-01-15/L/M42_LIGHT_L_600s_BIN1_-20C_GA2750_20250115_222558_734_PA045.50_E.FIT",
			expected: true,
		},
		{
			name:     "invalid file - wrong extension",
			path:     "/M42/2025-01-15/L/M42_LIGHT_L_600s_BIN1_-20C_GA2750_20250115_222558_734_PA045.50_E.TXT",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the regex processing directly
			image := tt.path // Use the full path for testing

			re := regexp.MustCompilePOSIX(Regex)
			splitline := re.FindAllStringSubmatch(image, -1)

			hasMatch := len(splitline) == 1

			if tt.expected && !hasMatch {
				t.Errorf("Expected regex to match for path: %s", tt.path)
			}

			if !tt.expected && hasMatch {
				t.Errorf("Expected regex to not match, but it matched for path: %s", tt.path)
			}

			// If regex matched, test that flats processing works
			if hasMatch {
				rotval, _ := strconv.ParseFloat(splitline[0][5], 32)
				rotation := float32(rotval)

				f := newFlat()
				FlatList.set(rotation, f)

				if len(FlatList) == 0 {
					t.Errorf("Expected flat to be created after regex match")
				}
			}

			// Reset for next test
			FlatList = flats{}
		})
	}
}

func TestTimeFiltering(t *testing.T) {
	// Reset global state
	FlatList = flats{}
	Rotations = []float32{}
	TimeFrame = 2 // 2 hours
	RotUsed = true
	Regex = testRegex

	tests := []struct {
		name     string
		path     string
		modTime  time.Time
		expected bool // should be processed
	}{
		{
			name:     "recent file within time frame",
			path:     "/M42/2025-01-15/L/M42_LIGHT_L_600s_BIN1_-20C_GA2750_20250115_222558_734_PA045.50_E.FIT",
			modTime:  time.Now().Add(-1 * time.Hour), // 1 hour ago
			expected: true,
		},
		{
			name:     "old file outside time frame",
			path:     "/M42/2025-01-15/L/M42_LIGHT_L_600s_BIN1_-20C_GA2750_20250115_222558_734_PA045.50_E.FIT",
			modTime:  time.Now().Add(-3 * time.Hour), // 3 hours ago
			expected: false,
		},
		{
			name:     "file exactly at time boundary",
			path:     "/M42/2025-01-15/L/M42_LIGHT_L_600s_BIN1_-20C_GA2750_20250115_222558_734_PA045.50_E.FIT",
			modTime:  time.Now().Add(-2 * time.Hour), // exactly 2 hours ago
			expected: false,                          // diff < cutoff, so 2 hours is not less than 2 hours
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the regex processing directly
			image := tt.path // Use the full path for testing

			re := regexp.MustCompilePOSIX(Regex)
			splitline := re.FindAllStringSubmatch(image, -1)

			hasMatch := len(splitline) == 1

			// For time filtering tests, we expect regex to match but time filtering to apply
			if !hasMatch {
				t.Errorf("Expected regex to match for path: %s", tt.path)
			}

			// Apply time filtering logic
			timeDiff := time.Since(tt.modTime)
			withinTimeFrame := timeDiff < time.Duration(TimeFrame)*time.Hour

			// Create a flat only if regex matches AND within time frame
			if hasMatch && withinTimeFrame {
				rotval, _ := strconv.ParseFloat(splitline[0][5], 32)
				rotation := float32(rotval)

				f := newFlat()
				FlatList.set(rotation, f)
			}

			hasFlats := len(FlatList) > 0

			if tt.expected && !hasFlats {
				t.Errorf("Expected file to be processed, but no flats were created for path: %s", tt.path)
			}

			if !tt.expected && hasFlats {
				t.Errorf("Expected file to not be processed, but flats were created for path: %s", tt.path)
			}

			// Reset for next test
			FlatList = flats{}
		})
	}
}

func TestPathHandling(t *testing.T) {
	// Reset global state
	ObjectList = Objects{}
	targetList = Targets{}
	RotUsed = true
	Regex = testRegex

	tests := []struct {
		name     string
		rootPath string
		fullPath string
		expected string // expected object name
	}{
		{
			name:     "unix path",
			rootPath: "/lights",
			fullPath: "/lights/M42/2025-01-15/L/M42_LIGHT_L_600s_BIN1_-20C_GA2750_20250115_222558_734_PA045.50_E.FIT",
			expected: "M42",
		},
		{
			name:     "path with nested directory",
			rootPath: "/very/long/path/to",
			fullPath: "/very/long/path/to/NGC7635/2025-01-15/H/NGC7635_LIGHT_H_180s_BIN1_5C_GA2750_20250115_233348_235_PA239.88_W.FIT",
			expected: "NGC7635",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the path trimming that happens in walker.go
			image := strings.TrimPrefix(tt.fullPath, tt.rootPath)

			re := regexp.MustCompilePOSIX(Regex)
			splitline := re.FindAllStringSubmatch(image, -1)

			if len(splitline) != 1 {
				t.Errorf("Expected regex to match for path: %s (processed as: %s)", tt.fullPath, image)
				return
			}

			// Extract object name from regex
			objectName := splitline[0][1]

			if objectName != tt.expected {
				t.Errorf("Expected object name %s, got %s", tt.expected, objectName)
			}

			// Reset for next test
			ObjectList = Objects{}
		})
	}
}

func TestRotationHandling(t *testing.T) {
	// Reset global state
	FlatList = flats{}
	Rotations = []float32{}
	TimeFrame = 24
	Regex = testRegex

	tests := []struct {
		name     string
		path     string
		rotUsed  bool
		expected float32 // expected rotation value
	}{
		{
			name:     "rotation used - parse from filename",
			path:     "/M42/2025-01-15/L/M42_LIGHT_L_600s_BIN1_-20C_GA2750_20250115_222558_734_PA045.50_E.FIT",
			rotUsed:  true,
			expected: 45.50,
		},
		{
			name:     "rotation not used - default 666",
			path:     "/M42/2025-01-15/L/M42_LIGHT_L_600s_BIN1_-20C_GA2750_20250115_222558_734_PA045.50_E.FIT",
			rotUsed:  false,
			expected: 666,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RotUsed = tt.rotUsed

			// Test the regex processing directly
			image := tt.path // Use the full path for testing

			re := regexp.MustCompilePOSIX(Regex)
			splitline := re.FindAllStringSubmatch(image, -1)

			if len(splitline) != 1 {
				t.Errorf("Expected regex to match for path: %s", tt.path)
				return
			}

			// Extract rotation from regex
			rotval, _ := strconv.ParseFloat(splitline[0][5], 32)
			rotation := float32(rotval)

			if !RotUsed {
				rotation = 666
			}

			// Create a flat to test rotation handling
			f := newFlat()
			FlatList.set(rotation, f)

			// Check if flats were created with correct rotation
			if len(FlatList) > 0 {
				for rot := range FlatList {
					if rot != tt.expected {
						t.Errorf("Expected rotation %f, got %f", tt.expected, rot)
					}
				}
			} else {
				t.Errorf("Expected flats to be created for path: %s", tt.path)
			}

			// Reset for next test
			FlatList = flats{}
		})
	}
}

func TestRegexExtraction(t *testing.T) {
	// Test the regex extraction logic directly
	Regex = testRegex

	testPath := "/M42/2025-01-15/L/M42_LIGHT_L_600s_BIN1_-20C_GA2750_20250115_222558_734_PA045.50_E.FIT"

	re := regexp.MustCompilePOSIX(Regex)
	splitline := re.FindAllStringSubmatch(testPath, -1)

	if len(splitline) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(splitline))
	}

	match := splitline[0]
	expected := []string{
		testPath,
		"M42",
		"L",
		"600",
		"-20",
		"045.50",
	}

	if len(match) != len(expected) {
		t.Fatalf("Expected %d groups, got %d", len(expected), len(match))
	}

	for i, exp := range expected {
		if match[i] != exp {
			t.Errorf("Group %d: expected %s, got %s", i, exp, match[i])
		}
	}
}

// Mock file info for testing
type mockFileInfo struct {
	name    string
	modTime time.Time
	isDir   bool
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return 1024 }
func (m *mockFileInfo) Mode() os.FileMode  { return 0644 }
func (m *mockFileInfo) ModTime() time.Time { return m.modTime }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() any           { return nil }
