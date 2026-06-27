package mac

import (
	"context"
	"os/exec"
	"senk-tui/managers"
	"sort"
	"strings"
	"time"
)

type BrewManager struct{}

func (m BrewManager) Name() string { return "Homebrew" }
func (m BrewManager) ManagerVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "brew", "--version").Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 0 {
		return lines[0], nil
	}
	return "", nil
}

func (m BrewManager) IsInstalled() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func (m BrewManager) Fetch() ([]managers.Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 1. Get all versions
	versionsCmd := exec.CommandContext(ctx, "brew", "list", "--versions")
	versionsOut, _ := versionsCmd.Output()
	allVersions := parseBrewVersionsMap(versionsOut)

	// 2. Get leaves (explicitly installed formulae)
	leavesCmd := exec.CommandContext(ctx, "brew", "leaves")
	leavesOut, _ := leavesCmd.Output()

	// 3. Get casks (explicitly installed Mac apps)
	casksCmd := exec.CommandContext(ctx, "brew", "list", "--cask")
	casksOut, _ := casksCmd.Output()

	return parseFilteredBrewOutput(leavesOut, casksOut, allVersions)
}

func parseBrewVersionsMap(data []byte) map[string]string {
	versions := make(map[string]string)
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			versions[parts[0]] = parts[1]
		}
	}
	return versions
}

func parseFilteredBrewOutput(leavesData, casksData []byte, versionsMap map[string]string) ([]managers.Dependency, error) {
	var deps []managers.Dependency
	seen := make(map[string]bool)

	// Add leaves
	leavesLines := strings.Split(string(leavesData), "\n")
	for _, line := range leavesLines {
		name := strings.TrimSpace(line)
		if name != "" && !seen[name] {
			version := versionsMap[name]
			if version == "" {
				version = "unknown"
			}
			deps = append(deps, managers.Dependency{Name: name, Version: version})
			seen[name] = true
		}
	}

	// Add casks
	caskLines := strings.Split(string(casksData), "\n")
	for _, line := range caskLines {
		name := strings.TrimSpace(line)
		if name != "" && !seen[name] {
			version := versionsMap[name]
			if version == "" {
				version = "unknown"
			}
			deps = append(deps, managers.Dependency{Name: name, Version: version})
			seen[name] = true
		}
	}

	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Name < deps[j].Name
	})

	if deps == nil {
		deps = []managers.Dependency{}
	}

	return deps, nil
}
