package managers
import (
	"reflect"
	"testing"
	"os"
)
func TestParseComposerOutput(t *testing.T) {
	validData, _ := os.ReadFile("testdata/composer_valid.json")
	tests := []struct {
		name     string
		input    []byte
		expected []Dependency
	}{
		{
			name: "Valid composer json",
			input: validData,
			expected: []Dependency{{Name: "laravel/installer", Version: "v5.1.0"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := parseComposerOutput(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
