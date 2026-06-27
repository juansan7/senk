package managers

import (
	"encoding/json"
	"sort"
)

type pipPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func parsePipOutput(data []byte) ([]Dependency, error) {
	var packages []pipPackage
	if err := json.Unmarshal(data, &packages); err != nil {
		return nil, err
	}

	var deps []Dependency
	for _, pkg := range packages {
		deps = append(deps, Dependency{
			Name:    pkg.Name,
			Version: pkg.Version,
		})
	}

	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Name < deps[j].Name
	})

	if deps == nil {
		deps = []Dependency{}
	}

	return deps, nil
}
