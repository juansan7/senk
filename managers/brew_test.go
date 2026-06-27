package managers

import (
	"reflect"
	"testing"
)

func TestParseFilteredBrewOutput(t *testing.T) {
	leavesData := []byte("jq\nnode\ngo\n")
	casksData := []byte("google-chrome\nspotify\n")
	versionsMap := map[string]string{
		"jq":            "1.7.1",
		"node":          "21.2.0",
		"go":            "1.21.6",
		"google-chrome": "120.0.0.0",
		"spotify":       "1.2.25",
		"openssl":       "3.0.0", // This shouldn't appear in the output because it's not a leaf/cask
	}

	expected := []Dependency{
		{Name: "go", Version: "1.21.6"},
		{Name: "google-chrome", Version: "120.0.0.0"},
		{Name: "jq", Version: "1.7.1"},
		{Name: "node", Version: "21.2.0"},
		{Name: "spotify", Version: "1.2.25"},
	}

	result, err := parseFilteredBrewOutput(leavesData, casksData, versionsMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
