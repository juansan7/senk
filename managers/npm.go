package managers

import (
	"encoding/json"
	"sort"
)

// npmOutput represents the exact JSON structure returned by `npm list -g --depth=0 --json`
type npmOutput struct {
	Dependencies map[string]struct {
		Version string `json:"version"`
	} `json:"dependencies"`
}

// parseNPMOutput is a pure function that takes JSON bytes and returns structured Dependencies.
func parseNPMOutput(data []byte) ([]Dependency, error) {
	var out npmOutput
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}

	var deps []Dependency
	for name, info := range out.Dependencies {
		deps = append(deps, Dependency{
			Name:    name,
			Version: info.Version,
		})
	}

	// Sort alphabetically by package name to ensure deterministic UI and easy testing.
	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Name < deps[j].Name
	})

	// If no dependencies were found, return an empty slice instead of nil
	if deps == nil {
		deps = []Dependency{}
	}

	return deps, nil
}
