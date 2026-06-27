package managers

import (
	"os"
	"reflect"
	"testing"
)

func TestParseBrewOutput(t *testing.T) {
	validData, err := os.ReadFile("testdata/brew_valid.txt")
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
			name:  "Valid brew output",
			input: validData,
			expected: []Dependency{
				{Name: "curl", Version: "8.4.0"},
				{Name: "git", Version: "2.42.0"},
				{Name: "go", Version: "1.21.4"},
				{Name: "jq", Version: "1.7.1"},
				{Name: "node", Version: "21.2.0"},
				{Name: "python@3.11", Version: "3.11.6"},
			},
			hasError: false,
		},
		{
			name:     "Empty Output",
			input:    []byte(``),
			expected: []Dependency{},
			hasError: false,
		},
		{
			name: "Malformed Line (ignored)",
			input: []byte(`justoneword
valid 1.0.0`),
			expected: []Dependency{
				{Name: "valid", Version: "1.0.0"},
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseBrewOutput(tt.input)

			if (err != nil) != tt.hasError {
				t.Fatalf("expected error: %v, got: %v", tt.hasError, err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
