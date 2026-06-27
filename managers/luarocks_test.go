package managers

import (
	"os"
	"reflect"
	"testing"
)

func TestParseLuarocksOutput(t *testing.T) {
	validData, _ := os.ReadFile("testdata/luarocks_valid.txt")
	tests := []struct {
		name     string
		input    []byte
		expected []Dependency
	}{
		{
			name:  "Valid luarocks output",
			input: validData,
			expected: []Dependency{
				{Name: "luasocket", Version: "3.0.0-1"},
				{Name: "penlight", Version: "1.13.1-1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := parseLuarocksOutput(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
