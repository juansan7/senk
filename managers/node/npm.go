package node

import (
	"context"
	"encoding/json"
	"os/exec"
	"senk-tui/managers"
	"sort"
	"strings"
	"time"
)

type NpmManager struct{}

func (m NpmManager) Name() string { return "Node.js (npm)" }

func (m NpmManager) IsInstalled() bool {
	_, err := exec.LookPath("npm")
	return err == nil
}

func (m NpmManager) Fetch() ([]managers.Dependency, error) {
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

func parseNPMOutput(data []byte) ([]managers.Dependency, error) {
	var out npmOutput
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}

	var deps []managers.Dependency
	for name, info := range out.Dependencies {
		deps = append(deps, managers.Dependency{
			Name:    name,
			Version: info.Version,
		})
	}

	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Name < deps[j].Name
	})

	if deps == nil {
		deps = []managers.Dependency{}
	}

	return deps, nil
}

func (m NpmManager) ManagerVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "npm", "--version").Output()
	if err != nil {
		return "", err
	}
	return "v" + strings.TrimSpace(string(out)), nil
}
