package managers
import (
	"reflect"
	"testing"
	"os"
)
func TestParsePubOutput(t *testing.T) {
	validData, _ := os.ReadFile("testdata/pub_valid.txt")
	tests := []struct {
		name     string
		input    []byte
		expected []Dependency
	}{
		{
			name: "Valid pub output",
			input: validData,
			expected: []Dependency{
				{Name: "fvm", Version: "2.4.1"},
				{Name: "melos", Version: "3.1.1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := parsePubOutput(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
