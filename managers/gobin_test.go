package managers

import (
	"reflect"
	"testing"
	"os"
)

func TestParseGoBinOutput(t *testing.T) {
	validData, err := os.ReadFile("testdata/gobin_valid.txt")
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
			name:  "Valid go version output",
			input: validData,
			expected: []Dependency{
				{Name: "air", Version: "v1.49.0"},
				{Name: "dlv", Version: "v1.21.0"},
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
			result, err := parseGoBinOutput(tt.input)

			if (err != nil) != tt.hasError {
				t.Fatalf("expected error: %v, got: %v", tt.hasError, err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
