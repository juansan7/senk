package managers

import (
	"context"
	"encoding/json"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type PipManager struct{}

func (m PipManager) Name() string { return "Python (pip)" }

func (m PipManager) IsInstalled() bool {
	_, err := exec.LookPath("pip")
	return err == nil
}

func (m PipManager) Fetch() ([]Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "pip", "list", "--not-required", "--format=json")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parsePipOutput(out)
}

type pipPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func parsePipOutput(data []byte) ([]Dependency, error) {
	var packages []pipPackage
	if len(data) > 0 {
		if err := json.Unmarshal(data, &packages); err != nil {
			return nil, err
		}
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

func (m PipManager) ManagerVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "pip", "--version").Output()
	if err != nil {
		return "", err
	}
	parts := strings.Fields(string(out))
	if len(parts) >= 2 {
		return "v" + parts[1], nil
	}
	return "", nil
}
