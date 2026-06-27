package managers
import (
	"reflect"
	"testing"
	"os"
)
func TestParseDotnetOutput(t *testing.T) {
	validData, _ := os.ReadFile("testdata/dotnet_valid.txt")
	tests := []struct {
		name     string
		input    []byte
		expected []Dependency
	}{
		{
			name: "Valid dotnet output",
			input: validData,
			expected: []Dependency{
				{Name: "dotnet-ef", Version: "7.0.5"},
				{Name: "dotnetsay", Version: "1.0.0"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := parseDotnetOutput(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
