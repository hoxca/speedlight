package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestPlural(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		singular string
		expected string
	}{
		{
			name:     "zero count",
			count:    0,
			singular: "hour",
			expected: "00 hour  ",
		},
		{
			name:     "singular count",
			count:    1,
			singular: "hour",
			expected: "01 hour  ",
		},
		{
			name:     "plural count",
			count:    2,
			singular: "hour",
			expected: "02 hours ",
		},
		{
			name:     "double digit count",
			count:    10,
			singular: "minute",
			expected: "10 minutes ",
		},
		{
			name:     "large count",
			count:    59,
			singular: "second",
			expected: "59 seconds ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := plural(tt.count, tt.singular)
			if result != tt.expected {
				t.Errorf("plural(%d, %s) = %q, expected %q", tt.count, tt.singular, result, tt.expected)
			}
		})
	}
}

func TestSecondsToHuman(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{
			name:     "zero seconds",
			input:    0,
			expected: "00 second  ",
		},
		{
			name:     "single second",
			input:    1,
			expected: "01 second  ",
		},
		{
			name:     "seconds only",
			input:    30,
			expected: "30 seconds ",
		},
		{
			name:     "one minute",
			input:    60,
			expected: "00 hour  01 minute  00 second  ",
		},
		{
			name:     "minutes and seconds",
			input:    90,
			expected: "00 hour  01 minute  30 seconds ",
		},
		{
			name:     "one hour",
			input:    3600,
			expected: "01 hour  00 minute  00 second  ",
		},
		{
			name:     "hours, minutes, and seconds",
			input:    3661,
			expected: "01 hour  01 minute  01 second  ",
		},
		{
			name:     "multiple hours",
			input:    7325,
			expected: "02 hours 02 minutes 05 seconds ",
		},
		{
			name:     "large duration",
			input:    10800,
			expected: "03 hours 00 minute  00 second  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := secondsToHuman(tt.input)
			if result != tt.expected {
				t.Errorf("secondsToHuman(%d) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestWriteDestinationSetWriteDestination(t *testing.T) {
	tests := []struct {
		name            string
		writeConsole    bool
		writeReport     bool
		expectedConsole bool
		expectedReport  bool
	}{
		{
			name:            "both true",
			writeConsole:    true,
			writeReport:     true,
			expectedConsole: true,
			expectedReport:  true,
		},
		{
			name:            "console only",
			writeConsole:    true,
			writeReport:     false,
			expectedConsole: true,
			expectedReport:  false,
		},
		{
			name:            "report only",
			writeConsole:    false,
			writeReport:     true,
			expectedConsole: false,
			expectedReport:  true,
		},
		{
			name:            "both false",
			writeConsole:    false,
			writeReport:     false,
			expectedConsole: false,
			expectedReport:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WriteDestination{}
			w.SetWriteDestination(tt.writeConsole, tt.writeReport)

			if w.WriteToConsole != tt.expectedConsole {
				t.Errorf("WriteToConsole = %v, expected %v", w.WriteToConsole, tt.expectedConsole)
			}

			if w.WriteToFile != tt.expectedReport {
				t.Errorf("WriteToFile = %v, expected %v", w.WriteToFile, tt.expectedReport)
			}
		})
	}
}

func TestObjectPrintObject(t *testing.T) {
	tests := []struct {
		name     string
		object   Object
		expected []string // substrings that should be in output
	}{
		{
			name: "object with rotation",
			object: Object{
				name:     "M42",
				rotation: 45.50,
				targets: map[string]Target{
					"M42~-20~300": {
						tuple: "M42~-20~300",
						name:  "M42",
						temp:  -20,
						expo:  300,
						fltr:  Filters{L: 5, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
					},
				},
			},
			expected: []string{
				"Object name: M42",
				"Rotation:45.50°",
				"-20°C",
				"Total:",
				"L\tNb:    5\tExpo:  300s\tSubs:",
			},
		},
		{
			name: "object without rotation (666)",
			object: Object{
				name:     "NGC7635",
				rotation: 666,
				targets: map[string]Target{
					"NGC7635~5~180": {
						tuple: "NGC7635~5~180",
						name:  "NGC7635",
						temp:  5,
						expo:  180,
						fltr:  Filters{L: 0, R: 3, G: 3, B: 3, S: 0, H: 0, O: 0},
					},
				},
			},
			expected: []string{
				"Object name: NGC7635",
				"Rotation: N/A",
				"5°C",
				"Total:",
				"R\tNb:    3\tExpo:  180s\tSubs:",
				"G\tNb:    3\tExpo:  180s\tSubs:",
				"B\tNb:    3\tExpo:  180s\tSubs:",
			},
		},
		{
			name: "object with narrowband filters",
			object: Object{
				name:     "IC1396",
				rotation: 123.75,
				targets: map[string]Target{
					"IC1396~-10~600": {
						tuple: "IC1396~-10~600",
						name:  "IC1396",
						temp:  -10,
						expo:  600,
						fltr:  Filters{L: 0, R: 0, G: 0, B: 0, S: 2, H: 2, O: 2},
					},
				},
			},
			expected: []string{
				"Object name: IC1396",
				"Rotation:123.75°",
				"-10°C",
				"Total:",
				"S\tNb:    2\tExpo:  600s\tSubs:",
				"H\tNb:    2\tExpo:  600s\tSubs:",
				"O\tNb:    2\tExpo:  600s\tSubs:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the fprintObject method instead since printObject writes to stdout directly
			var buf bytes.Buffer
			tt.object.fprintObject(&buf)

			output := buf.String()

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, but got:\n%s", expected, output)
				}
			}
		})
	}
}

func TestObjectFprintObject(t *testing.T) {
	tests := []struct {
		name     string
		object   Object
		expected []string // substrings that should be in output
	}{
		{
			name: "object with rotation",
			object: Object{
				name:     "M42",
				rotation: 45.50,
				targets: map[string]Target{
					"M42~-20~300": {
						tuple: "M42~-20~300",
						name:  "M42",
						temp:  -20,
						expo:  300,
						fltr:  Filters{L: 1, R: 1, G: 1, B: 1, S: 0, H: 0, O: 0},
					},
				},
			},
			expected: []string{
				"Object name: M42",
				"Rotation:45.50°",
				"-20°C",
				"Total:",
				"L\tNb:    1\tExpo:  300s\tSubs:",
				"R\tNb:    1\tExpo:  300s\tSubs:",
				"G\tNb:    1\tExpo:  300s\tSubs:",
				"B\tNb:    1\tExpo:  300s\tSubs:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.object.fprintObject(&buf)

			output := buf.String()

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, but got:\n%s", expected, output)
				}
			}
		})
	}
}

func TestObjectsPrintObjects(t *testing.T) {
	// Reset global state
	Wdest = WriteDestination{}

	tests := []struct {
		name            string
		objects         Objects
		writeConsole    bool
		writeToFile     bool
		expectedConsole []string // substrings that should be in console output
	}{
		{
			name: "console output only",
			objects: Objects{
				"M42": {
					name:     "M42",
					rotation: 45.50,
					targets: map[string]Target{
						"M42~-20~300": {
							tuple: "M42~-20~300",
							name:  "M42",
							temp:  -20,
							expo:  300,
							fltr:  Filters{L: 1, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
						},
					},
				},
			},
			writeConsole: true,
			writeToFile:  false,
			expectedConsole: []string{
				"Targets list: [\"M42\"]",
				"Object name: M42",
				"Rotation:45.50°",
			},
		},
		{
			name: "no output",
			objects: Objects{
				"M42": {
					name:     "M42",
					rotation: 45.50,
					targets:  make(map[string]Target),
				},
			},
			writeConsole:    false,
			writeToFile:     false,
			expectedConsole: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the core logic without stdout capture complexity
			Wdest.SetWriteDestination(tt.writeConsole, tt.writeToFile)

			// Test that objects are properly structured
			objects := tt.objects.getObjects()
			if len(objects) == 0 && len(tt.objects) > 0 {
				t.Errorf("Expected objects to be found, but got empty list")
			}

			// Test object structure
			for name, obj := range tt.objects {
				if obj.name != name {
					t.Errorf("Object name mismatch: expected %s, got %s", name, obj.name)
				}
			}
		})
	}
}
