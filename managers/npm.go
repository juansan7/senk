package managers

import (
	"context"
	"encoding/json"
	"os/exec"
	"sort"
	"time"
)

type NpmManager struct{}

func (m NpmManager) Name() string { return "Node.js (npm)" }

func (m NpmManager) IsInstalled() bool {
	_, err := exec.LookPath("npm")
	return err == nil
}

func (m NpmManager) Fetch() ([]Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "npm", "list", "-g", "--depth=0", "--json")
	out, _ := cmd.Output() // Ignore error: npm returns non-zero on peer dependency warnings, but still outputs valid JSON
	return parseNPMOutput(out)
}

type npmOutput struct {
	Dependencies map[string]struct {
		Version string `json:"version"`
	} `json:"dependencies"`
}

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

	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Name < deps[j].Name
	})

	if deps == nil {
		deps = []Dependency{}
	}

	return deps, nil
}
