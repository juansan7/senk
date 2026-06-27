package managers

import (
	"sort"
	"strings"
)

func parseGoBinOutput(data []byte) ([]Dependency, error) {
	var deps []Dependency
	lines := strings.Split(string(data), "\n")

	var currentName string
	var currentVersion string

	for _, line := range lines {
		// Output format:
		// /path/to/bin: go1.x.y
		// \tpath\tgithub.com/user/repo/cmd/name
		// \tmod\tgithub.com/user/repo\tv1.2.3\th1:...
		if !strings.HasPrefix(line, "\t") {
			// This is the binary path line. We can extract the name from the path.
			// e.g. /Users/admin/go/bin/dlv: go1.21.0
			parts := strings.Split(line, ":")
			if len(parts) > 0 {
				pathParts := strings.Split(parts[0], "/")
				if len(pathParts) > 0 {
					name := pathParts[len(pathParts)-1]
					currentName = strings.TrimSpace(name)
					currentVersion = "(devel)" // fallback
				}
			}
		} else if strings.HasPrefix(line, "\tmod\t") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				currentVersion = parts[2]
			}
			if currentName != "" {
				deps = append(deps, Dependency{
					Name:    currentName,
					Version: currentVersion,
				})
				currentName = ""
			}
		} else if strings.HasPrefix(line, "\tpath\t") {
			// fallback in case 'mod' line is missing but 'path' is present.
			// Actually we usually wait for 'mod' to get version. 
			// If we only have 'path', it might be a local build, but we'll wait for next lines.
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
