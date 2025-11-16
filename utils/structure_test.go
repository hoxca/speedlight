package utils

import (
	"reflect"
	"testing"
)

func TestFiltersIterateFilter(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		initial  Filters
		expected Filters
	}{
		{
			name:     "increment L filter",
			filter:   "L",
			initial:  Filters{},
			expected: Filters{L: 1},
		},
		{
			name:     "increment R filter",
			filter:   "R",
			initial:  Filters{L: 2},
			expected: Filters{L: 2, R: 1},
		},
		{
			name:     "increment G filter",
			filter:   "G",
			initial:  Filters{},
			expected: Filters{G: 1},
		},
		{
			name:     "increment B filter",
			filter:   "B",
			initial:  Filters{},
			expected: Filters{B: 1},
		},
		{
			name:     "increment S filter",
			filter:   "S",
			initial:  Filters{},
			expected: Filters{S: 1},
		},
		{
			name:     "increment H filter",
			filter:   "H",
			initial:  Filters{},
			expected: Filters{H: 1},
		},
		{
			name:     "increment O filter",
			filter:   "O",
			initial:  Filters{},
			expected: Filters{O: 1},
		},
		{
			name:     "unknown filter no change",
			filter:   "X",
			initial:  createTestFilters(),
			expected: createTestFilters(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.initial
			f.iterateFilter(tt.filter)
			assertFiltersEqual(t, f, tt.expected)
		})
	}
}

func TestNewTarget(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		expected *Target
	}{
		{
			name:   "valid target string",
			target: "M42~-20~300",
			expected: &Target{
				tuple: "M42~-20~300",
				name:  "M42",
				temp:  -20,
				expo:  300,
				fltr:  Filters{},
			},
		},
		{
			name:   "target with positive temperature",
			target: "NGC7635~5~180",
			expected: &Target{
				tuple: "NGC7635~5~180",
				name:  "NGC7635",
				temp:  5,
				expo:  180,
				fltr:  Filters{},
			},
		},
		{
			name:   "malformed target string",
			target: "M42~invalid~300",
			expected: &Target{
				tuple: "M42~invalid~300",
				name:  "M42",
				temp:  0,
				expo:  300,
				fltr:  Filters{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newTarget(tt.target)
			assertTargetEqual(t, result, tt.expected)
		})
	}
}

func TestNewObject(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		rotation float32
		expected *Object
	}{
		{
			name:     "valid object creation",
			target:   "M42~-20~300",
			rotation: 45.5,
			expected: &Object{
				name:     "M42",
				rotation: 45.5,
				targets:  make(map[string]Target),
			},
		},
		{
			name:     "object with special rotation",
			target:   "NGC7635~5~180",
			rotation: 666,
			expected: &Object{
				name:     "NGC7635",
				rotation: 666,
				targets:  make(map[string]Target),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newObject(tt.target, tt.rotation)
			assertObjectEqual(t, result, tt.expected)
		})
	}
}

func TestObjectsExist(t *testing.T) {
	tests := []struct {
		name     string
		objects  Objects
		key      string
		expected bool
	}{
		{
			name:     "existing key",
			objects:  Objects{"M42": *createTestObject("M42", 45.5)},
			key:      "M42",
			expected: true,
		},
		{
			name:     "non-existing key",
			objects:  Objects{"M42": *createTestObject("M42", 45.5)},
			key:      "NGC7635",
			expected: false,
		},
		{
			name:     "empty objects",
			objects:  Objects{},
			key:      "M42",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.objects.exist(tt.key)
			if result != tt.expected {
				t.Errorf("Objects.exist() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestTargetsExist(t *testing.T) {
	tests := []struct {
		name     string
		targets  Targets
		key      string
		expected bool
	}{
		{
			name:     "existing key",
			targets:  Targets{"M42~-20~300": *createTestTarget("M42~-20~300")},
			key:      "M42~-20~300",
			expected: true,
		},
		{
			name:     "non-existing key",
			targets:  Targets{"M42~-20~300": *createTestTarget("M42~-20~300")},
			key:      "M42~-15~300",
			expected: false,
		},
		{
			name:     "empty targets",
			targets:  Targets{},
			key:      "M42~-20~300",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.targets.exist(tt.key)
			if result != tt.expected {
				t.Errorf("Targets.exist() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestObjectsSet(t *testing.T) {
	objects := Objects{}
	obj := createTestObject("M42", 45.5)

	objects.set("M42", obj)

	if len(objects) != 1 {
		t.Errorf("Objects.set() expected 1 object, got %d", len(objects))
	}

	if _, exists := objects["M42"]; !exists {
		t.Error("Objects.set() key not found in map")
	}
}

func TestObjectSetTarget(t *testing.T) {
	obj := newObject("M42", 45.5)
	target := createTestTarget("M42~-20~300")

	obj.setTarget("M42~-20~300", target)

	if len(obj.targets) != 1 {
		t.Errorf("Object.setTarget() expected 1 target, got %d", len(obj.targets))
	}

	if _, exists := obj.targets["M42~-20~300"]; !exists {
		t.Error("Object.setTarget() key not found in targets map")
	}
}

func TestObjectsGetObjects(t *testing.T) {
	tests := []struct {
		name     string
		objects  Objects
		expected []string
	}{
		{
			name: "multiple objects unsorted",
			objects: Objects{
				"Z": *createTestObject("Z", 0),
				"A": *createTestObject("A", 0),
				"M": *createTestObject("M", 0),
			},
			expected: []string{"A", "M", "Z"},
		},
		{
			name:     "empty objects",
			objects:  Objects{},
			expected: []string{},
		},
		{
			name: "single object",
			objects: Objects{
				"M42": *createTestObject("M42", 45.5),
			},
			expected: []string{"M42"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.objects.getObjects()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Objects.getObjects() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestTargetsGetTarget(t *testing.T) {
	target := createTestTarget("M42~-20~300")
	targets := Targets{
		"M42~-20~300":   *target,
		"NGC7635~5~180": *createTestTarget("NGC7635~5~180"),
	}

	tests := []struct {
		name     string
		targets  Targets
		key      string
		expected *Target
	}{
		{
			name:     "existing target",
			targets:  targets,
			key:      "M42~-20~300",
			expected: target,
		},
		{
			name:     "non-existing target",
			targets:  targets,
			key:      "M31~-10~600",
			expected: nil,
		},
		{
			name:     "empty targets",
			targets:  Targets{},
			key:      "M42~-20~300",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.targets.getTarget(tt.key)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Targets.getTarget() = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}

func TestObjectsGetObject(t *testing.T) {
	obj := createTestObject("M42", 45.5)
	objects := Objects{
		"M42":     *obj,
		"NGC7635": *createTestObject("NGC7635", 123.0),
	}

	tests := []struct {
		name     string
		objects  Objects
		key      string
		expected *Object
	}{
		{
			name:     "existing object",
			objects:  objects,
			key:      "M42",
			expected: obj,
		},
		{
			name:     "non-existing object",
			objects:  objects,
			key:      "M31",
			expected: nil,
		},
		{
			name:     "empty objects",
			objects:  Objects{},
			key:      "M42",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.objects.getObject(tt.key)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Objects.getObject() = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}

func TestAddTarget(t *testing.T) {
	setupTestEnvironment(t)

	tests := []struct {
		name     string
		target   string
		rotation float32
		expected string
	}{
		{
			name:     "add new target to new object",
			target:   "M42~-20~300",
			rotation: 45.5,
			expected: "M42",
		},
		{
			name:     "add new target to existing object",
			target:   "M42~-15~300",
			rotation: 45.5,
			expected: "M42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addTarget(tt.target, tt.rotation)
			if result == nil {
				t.Error("addTarget() returned nil")
			}
			if result.name != tt.expected {
				t.Errorf("addTarget() object name = %v, expected %v", result.name, tt.expected)
			}
		})
	}
}

func TestObjectAddFilter(t *testing.T) {
	obj := newObject("M42", 45.5)
	target := newTarget("M42~-20~300")
	obj.setTarget("M42~-20~300", target)

	tests := []struct {
		name     string
		filter   string
		expected Filters
	}{
		{
			name:   "add L filter",
			filter: "L",
			expected: Filters{
				L: 1, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0,
			},
		},
		{
			name:   "add R filter",
			filter: "R",
			expected: Filters{
				L: 1, R: 1, G: 0, B: 0, S: 0, H: 0, O: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj.addFilter("M42~-20~300", tt.filter)
			updatedTarget := obj.targets["M42~-20~300"]
			assertFiltersEqual(t, updatedTarget.fltr, tt.expected)
		})
	}
}
