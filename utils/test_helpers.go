package utils

import (
	"testing"
)

// resetGlobalState resets all global variables used in the utils package
// This should be called at the beginning of each test that depends on global state
func resetGlobalState() {
	ObjectList = Objects{}
	targetList = Targets{}
	FlatList = flats{}
	Rotations = []float32{}
	RotUsed = false
	TimeFrame = 14 // default value
	Regex = ""     // will be set by individual tests
	Wdest = WriteDestination{}

	// Reset walker globals
	i = true
	rootp = ""
}

// createTestObject creates a test Object with sample data
func createTestObject(name string, rotation float32) *Object {
	return &Object{
		name:     name,
		rotation: rotation,
		targets: map[string]Target{
			name + "~-20~300": {
				tuple: name + "~-20~300",
				name:  name,
				temp:  -20,
				expo:  300,
				fltr:  Filters{L: 5, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
			},
		},
	}
}

// createTestTarget creates a test Target with sample data
func createTestTarget(tuple string) *Target {
	return &Target{
		tuple: tuple,
		name:  "M42",
		temp:  -20,
		expo:  300,
		fltr:  Filters{L: 1, R: 1, G: 1, B: 1, S: 0, H: 0, O: 0},
	}
}

// createTestFilters creates test Filters with sample data
func createTestFilters() Filters {
	return Filters{
		L: 5,
		R: 3,
		G: 2,
		B: 1,
		S: 0,
		H: 0,
		O: 0,
	}
}

// assertFiltersEqual compares two Filters structs for testing
func assertFiltersEqual(t *testing.T, got, expected Filters) {
	if got.L != expected.L || got.R != expected.R || got.G != expected.G ||
		got.B != expected.B || got.S != expected.S || got.H != expected.H || got.O != expected.O {
		t.Errorf("Filters mismatch: got %+v, expected %+v", got, expected)
	}
}

// assertObjectEqual compares two Object structs for testing
func assertObjectEqual(t *testing.T, got, expected *Object) {
	if got.name != expected.name {
		t.Errorf("Object name mismatch: got %s, expected %s", got.name, expected.name)
	}
	if got.rotation != expected.rotation {
		t.Errorf("Object rotation mismatch: got %f, expected %f", got.rotation, expected.rotation)
	}
	if len(got.targets) != len(expected.targets) {
		t.Errorf("Object targets count mismatch: got %d, expected %d", len(got.targets), len(expected.targets))
	}
}

// assertTargetEqual compares two Target structs for testing
func assertTargetEqual(t *testing.T, got, expected *Target) {
	if got.tuple != expected.tuple {
		t.Errorf("Target tuple mismatch: got %s, expected %s", got.tuple, expected.tuple)
	}
	if got.name != expected.name {
		t.Errorf("Target name mismatch: got %s, expected %s", got.name, expected.name)
	}
	if got.temp != expected.temp {
		t.Errorf("Target temp mismatch: got %d, expected %d", got.temp, expected.temp)
	}
	if got.expo != expected.expo {
		t.Errorf("Target expo mismatch: got %d, expected %d", got.expo, expected.expo)
	}
	assertFiltersEqual(t, got.fltr, expected.fltr)
}

// setupTestEnvironment initializes a clean test environment
func setupTestEnvironment(t *testing.T) {
	resetGlobalState()

	// Set up default test values
	TimeFrame = 24
	Regex = `/(.*)/[[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2}/(.*)/.*_LIGHT_[LRGBSHO]*_([[:digit:]]*)s_BIN1_(.*)C_GA.*_[[:digit:]]{8}_[[:digit:]]{6}_[[:digit:]]{3}_PA([[:digit:]]{3}\.[[:digit:]]{2})_[EW]\.FIT`
	RotUsed = true
}

// cleanupTestEnvironment cleans up after a test
func cleanupTestEnvironment(t *testing.T) {
	resetGlobalState()
}
