package python

import (
	"os"
	"reflect"
	"senk-tui/managers"
	"testing"
)

func TestParsePipOutput(t *testing.T) {
	validData, err := os.ReadFile("testdata/pip_valid.json")
	if err != nil {
		t.Fatalf("Failed to read testdata: %v", err)
	}

	tests := []struct {
		name     string
		input    []byte
		expected []managers.Dependency
		hasError bool
	}{
		{
			name:  "Valid JSON output",
			input: validData,
			expected: []managers.Dependency{
				{Name: "certifi", Version: "2023.11.17"},
				{Name: "charset-normalizer", Version: "3.3.2"},
				{Name: "idna", Version: "3.6"},
				{Name: "requests", Version: "2.31.0"},
				{Name: "urllib3", Version: "2.1.0"},
			},
			hasError: false,
		},
		{
			name:     "Invalid JSON",
			input:    []byte(`[{"name": "broken"`),
			expected: nil,
			hasError: true,
		},
		{
			name:     "Empty Output",
			input:    []byte(`[]`),
			expected: []managers.Dependency{},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parsePipOutput(tt.input)

			if (err != nil) != tt.hasError {
				t.Fatalf("expected error: %v, got: %v", tt.hasError, err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
