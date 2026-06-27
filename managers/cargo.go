package managers

import (
	"regexp"
	"sort"
	"strings"
)

func parseCargoOutput(data []byte) ([]Dependency, error) {
	var deps []Dependency
	lines := strings.Split(string(data), "\n")

	// Cargo install --list format:
	// package-name v1.2.3:
	//     binary-name
	re := regexp.MustCompile(`^([^ ]+) v([^:]+):`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			deps = append(deps, Dependency{
				Name:    matches[1],
				Version: "v" + matches[2],
			})
		}
	}

	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Name < deps[j].Name
	})

	if deps == nil {
		deps = []Dependency{}
	}

	return deps, nil
}
