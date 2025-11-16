package utils

import (
	"reflect"
	"testing"
)

func TestAddFlats(t *testing.T) {
	// Reset global state before test
	FlatList = flats{}
	Rotations = []float32{}

	tests := []struct {
		name     string
		rotation float32
		filter   string
		expected map[float32]Filters
	}{
		{
			name:     "add first flat with L filter",
			rotation: 45.5,
			filter:   "L",
			expected: map[float32]Filters{
				45.5: {L: 1, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
			},
		},
		{
			name:     "add second flat with same rotation different filter",
			rotation: 45.5,
			filter:   "R",
			expected: map[float32]Filters{
				45.5: {L: 1, R: 1, G: 0, B: 0, S: 0, H: 0, O: 0},
			},
		},
		{
			name:     "add flat with different rotation",
			rotation: 123.75,
			filter:   "H",
			expected: map[float32]Filters{
				45.5:   {L: 1, R: 1, G: 0, B: 0, S: 0, H: 0, O: 0},
				123.75: {L: 0, R: 0, G: 0, B: 0, S: 0, H: 1, O: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addFlats(tt.rotation, tt.filter)

			if !reflect.DeepEqual(map[float32]Filters(FlatList), tt.expected) {
				t.Errorf("addFlats() = %+v, expected %+v", FlatList, tt.expected)
			}
		})
	}
}

func TestFlatsSet(t *testing.T) {
	f := flats{}
	filter := Filters{L: 5, R: 3, G: 2, B: 1, S: 0, H: 0, O: 0}

	f.set(45.5, &filter)

	if len(f) != 1 {
		t.Errorf("flats.set() expected 1 item, got %d", len(f))
	}

	if _, exists := f[45.5]; !exists {
		t.Error("flats.set() key not found in map")
	}

	result := f[45.5]
	if !reflect.DeepEqual(result, filter) {
		t.Errorf("flats.set() = %+v, expected %+v", result, filter)
	}
}

func TestFlatsExist(t *testing.T) {
	tests := []struct {
		name     string
		flats    flats
		key      float32
		expected bool
	}{
		{
			name:     "existing key",
			flats:    flats{45.5: {L: 1}},
			key:      45.5,
			expected: true,
		},
		{
			name:     "non-existing key",
			flats:    flats{45.5: {L: 1}},
			key:      123.75,
			expected: false,
		},
		{
			name:     "empty flats",
			flats:    flats{},
			key:      45.5,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.flats.exist(tt.key)
			if result != tt.expected {
				t.Errorf("flats.exist() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestNewFlat(t *testing.T) {
	result := newFlat()

	expected := Filters{
		L: 0, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0,
	}

	if !reflect.DeepEqual(*result, expected) {
		t.Errorf("newFlat() = %+v, expected %+v", *result, expected)
	}
}

func TestFlatsFiltersIterateFilter(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		initial  Filters
		expected Filters
	}{
		{
			name:     "increment L filter",
			filter:   "L",
			initial:  Filters{L: 0, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
			expected: Filters{L: 1, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
		},
		{
			name:     "increment R filter",
			filter:   "R",
			initial:  Filters{L: 2, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
			expected: Filters{L: 2, R: 1, G: 0, B: 0, S: 0, H: 0, O: 0},
		},
		{
			name:     "increment G filter",
			filter:   "G",
			initial:  Filters{L: 0, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
			expected: Filters{L: 0, R: 0, G: 1, B: 0, S: 0, H: 0, O: 0},
		},
		{
			name:     "increment B filter",
			filter:   "B",
			initial:  Filters{L: 0, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
			expected: Filters{L: 0, R: 0, G: 0, B: 1, S: 0, H: 0, O: 0},
		},
		{
			name:     "increment S filter",
			filter:   "S",
			initial:  Filters{L: 0, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
			expected: Filters{L: 0, R: 0, G: 0, B: 0, S: 1, H: 0, O: 0},
		},
		{
			name:     "increment H filter",
			filter:   "H",
			initial:  Filters{L: 0, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
			expected: Filters{L: 0, R: 0, G: 0, B: 0, S: 0, H: 1, O: 0},
		},
		{
			name:     "increment O filter",
			filter:   "O",
			initial:  Filters{L: 0, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
			expected: Filters{L: 0, R: 0, G: 0, B: 0, S: 0, H: 0, O: 1},
		},
		{
			name:     "unknown filter no change",
			filter:   "X",
			initial:  Filters{L: 1, R: 2, G: 3, B: 4, S: 5, H: 6, O: 7},
			expected: Filters{L: 1, R: 2, G: 3, B: 4, S: 5, H: 6, O: 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.initial
			f.iterateFilter(tt.filter)
			if !reflect.DeepEqual(f, tt.expected) {
				t.Errorf("iterateFilter() = %+v, expected %+v", f, tt.expected)
			}
		})
	}
}

func TestFiltersString(t *testing.T) {
	tests := []struct {
		name     string
		filters  Filters
		expected string
	}{
		{
			name:     "empty filters",
			filters:  Filters{L: 0, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
			expected: "",
		},
		{
			name:     "single L filter",
			filters:  Filters{L: 1, R: 0, G: 0, B: 0, S: 0, H: 0, O: 0},
			expected: "L",
		},
		{
			name:     "single R filter",
			filters:  Filters{L: 0, R: 1, G: 0, B: 0, S: 0, H: 0, O: 0},
			expected: "R",
		},
		{
			name:     "RGB filters",
			filters:  Filters{L: 1, R: 1, G: 1, B: 1, S: 0, H: 0, O: 0},
			expected: "LRGB",
		},
		{
			name:     "narrowband filters",
			filters:  Filters{L: 0, R: 0, G: 0, B: 0, S: 1, H: 1, O: 1},
			expected: "SHO",
		},
		{
			name:     "all filters",
			filters:  Filters{L: 1, R: 1, G: 1, B: 1, S: 1, H: 1, O: 1},
			expected: "LRGBSHO",
		},
		{
			name:     "mixed filters",
			filters:  Filters{L: 2, R: 0, G: 1, B: 0, S: 3, H: 0, O: 1},
			expected: "LGSO",
		},
		{
			name:     "only L and O filters",
			filters:  Filters{L: 5, R: 0, G: 0, B: 0, S: 0, H: 0, O: 2},
			expected: "LO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filters.String()
			if result != tt.expected {
				t.Errorf("Filters.String() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestAppendIfMissing(t *testing.T) {
	tests := []struct {
		name     string
		slice    []float32
		item     float32
		expected []float32
	}{
		{
			name:     "append to empty slice",
			slice:    []float32{},
			item:     45.5,
			expected: []float32{45.5},
		},
		{
			name:     "append new item",
			slice:    []float32{45.5},
			item:     123.75,
			expected: []float32{45.5, 123.75},
		},
		{
			name:     "item already exists",
			slice:    []float32{45.5, 123.75},
			item:     45.5,
			expected: []float32{45.5, 123.75},
		},
		{
			name:     "multiple items with duplicate",
			slice:    []float32{0.0, 45.5, 90.0},
			item:     90.0,
			expected: []float32{0.0, 45.5, 90.0},
		},
		{
			name:     "append special rotation 666",
			slice:    []float32{45.5},
			item:     666,
			expected: []float32{45.5, 666},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appendIfMissing(tt.slice, tt.item)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("appendIfMissing() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestRotationsManagement(t *testing.T) {
	// Reset global state before test
	FlatList = flats{}
	Rotations = []float32{}

	// Test that Rotations slice is properly managed when adding flats
	addFlats(45.5, "L")
	addFlats(123.75, "H")
	addFlats(45.5, "R") // same rotation, different filter

	expectedRotations := []float32{45.5, 123.75}

	if len(Rotations) != len(expectedRotations) {
		t.Errorf("Expected %d rotations, got %d", len(expectedRotations), len(Rotations))
	}

	for i, expected := range expectedRotations {
		if i >= len(Rotations) || Rotations[i] != expected {
			t.Errorf("Rotation %d: expected %f, got %f", i, expected, Rotations[i])
		}
	}
}

func TestFlatListIntegration(t *testing.T) {
	// Reset global state before test
	FlatList = flats{}
	Rotations = []float32{}

	// Add multiple flats with different rotations and filters
	addFlats(45.5, "L")
	addFlats(45.5, "R")
	addFlats(45.5, "G")
	addFlats(123.75, "H")
	addFlats(123.75, "S")
	addFlats(123.75, "O")

	// Check that we have the correct number of rotations
	if len(Rotations) != 2 {
		t.Errorf("Expected 2 rotations, got %d", len(Rotations))
	}

	// Check that each rotation has the correct filters
	filters45_5 := FlatList[45.5]
	expected45_5 := Filters{L: 1, R: 1, G: 1, B: 0, S: 0, H: 0, O: 0}
	if !reflect.DeepEqual(filters45_5, expected45_5) {
		t.Errorf("Rotation 45.5 filters = %+v, expected %+v", filters45_5, expected45_5)
	}

	filters123_75 := FlatList[123.75]
	expected123_75 := Filters{L: 0, R: 0, G: 0, B: 0, S: 1, H: 1, O: 1}
	if !reflect.DeepEqual(filters123_75, expected123_75) {
		t.Errorf("Rotation 123.75 filters = %+v, expected %+v", filters123_75, expected123_75)
	}
}
