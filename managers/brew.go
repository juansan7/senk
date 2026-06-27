package managers

import (
	"context"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type BrewManager struct{}

func (m BrewManager) Name() string { return "Homebrew" }

func (m BrewManager) IsInstalled() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func (m BrewManager) Fetch() ([]Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // brew can be slow
	defer cancel()
	cmd := exec.CommandContext(ctx, "brew", "list", "--versions")
	out, err := cmd.Output()
	
	// Homebrew sometimes returns exit status 1 (e.g. Ruby errors with casks)
	// but still prints the valid list of packages to stdout.
	// If we have output, we should try to parse it regardless of the error.
	if len(out) > 0 {
		return parseBrewOutput(out)
	}
	
	if err != nil {
		return nil, err
	}
	return parseBrewOutput(out)
}

func parseBrewOutput(data []byte) ([]Dependency, error) {
	var deps []Dependency
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			deps = append(deps, Dependency{
				Name:    parts[0],
				Version: parts[1],
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
