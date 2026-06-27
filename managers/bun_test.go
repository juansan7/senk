package managers
import (
	"reflect"
	"testing"
	"os"
)
func TestParseBunOutput(t *testing.T) {
	validData, _ := os.ReadFile("testdata/bun_valid.txt")
	tests := []struct {
		name     string
		input    []byte
		expected []Dependency
	}{
		{
			name: "Valid bun output",
			input: validData,
			expected: []Dependency{
				{Name: "cowsay", Version: "1.5.0"},
				{Name: "typescript", Version: "5.3.3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := parseBunOutput(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
