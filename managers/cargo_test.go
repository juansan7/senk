package managers

import (
	"reflect"
	"testing"
	"os"
)

func TestParseCargoOutput(t *testing.T) {
	validData, err := os.ReadFile("testdata/cargo_valid.txt")
	if err != nil {
		t.Fatalf("Failed to read testdata: %v", err)
	}

	tests := []struct {
		name     string
		input    []byte
		expected []Dependency
		hasError bool
	}{
		{
			name:  "Valid cargo output",
			input: validData,
			expected: []Dependency{
				{Name: "bat", Version: "v0.23.0"},
				{Name: "cargo-expand", Version: "v1.0.15"},
				{Name: "ripgrep", Version: "v13.0.0"},
			},
			hasError: false,
		},
		{
			name:     "Empty Output",
			input:    []byte(``),
			expected: []Dependency{},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseCargoOutput(tt.input)

			if (err != nil) != tt.hasError {
				t.Fatalf("expected error: %v, got: %v", tt.hasError, err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
