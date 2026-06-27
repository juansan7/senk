package managers

import (
	"reflect"
	"testing"
	"os"
)

func TestParseNPMOutput(t *testing.T) {
	validData, err := os.ReadFile("testdata/npm_valid.json")
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
			name:  "Valid JSON output",
			input: validData,
			expected: []Dependency{
				{Name: "npm", Version: "10.2.4"},
				{Name: "typescript", Version: "5.3.3"},
				{Name: "yarn", Version: "1.22.19"},
			},
			hasError: false,
		},
		{
			name:     "Invalid JSON",
			input:    []byte(`{"dependencies": { "broken": `),
			expected: nil,
			hasError: true,
		},
		{
			name:     "Empty Output",
			input:    []byte(`{}`),
			expected: []Dependency{},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseNPMOutput(tt.input)

			if (err != nil) != tt.hasError {
				t.Fatalf("expected error: %v, got: %v", tt.hasError, err)
			}

			// For unordered slices, we should compare lengths and elements individually.
			// But since we will sort the output alphabetically in our parser, DeepEqual is safe!
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
