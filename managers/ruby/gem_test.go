package ruby

import (
	"os"
	"reflect"
	"senk-tui/managers"
	"testing"
)

func TestParseGemOutput(t *testing.T) {
	validData, _ := os.ReadFile("testdata/gem_valid.txt")
	tests := []struct {
		name     string
		input    []byte
		expected []managers.Dependency
	}{
		{
			name:  "Valid gem output filtering defaults",
			input: validData,
			expected: []managers.Dependency{
				{Name: "CFPropertyList", Version: "2.3.6"},
				{Name: "actionmailer", Version: "7.1.2, 7.0.8"},
				{Name: "bundler", Version: "2.4.22"},
			},
		},
		{
			name:     "Empty",
			input:    []byte(""),
			expected: []managers.Dependency{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := parseGemOutput(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
