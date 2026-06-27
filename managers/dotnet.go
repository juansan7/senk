package managers

import (
	"context"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type DotnetManager struct{}

func (m DotnetManager) Name() string { return ".NET (dotnet)" }
func (m DotnetManager) IsInstalled() bool {
	_, err := exec.LookPath("dotnet")
	return err == nil
}
func (m DotnetManager) Fetch() ([]Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "dotnet", "tool", "list", "-g").Output()
	if len(out) > 0 {
		return parseDotnetOutput(out)
	}
	return nil, err
}

func parseDotnetOutput(data []byte) ([]Dependency, error) {
	var deps []Dependency
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if i < 2 {
			continue
		} // Skip headers
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			deps = append(deps, Dependency{Name: parts[0], Version: parts[1]})
		}
	}
	sort.Slice(deps, func(i, j int) bool { return deps[i].Name < deps[j].Name })
	if deps == nil {
		deps = []Dependency{}
	}
	return deps, nil
}

func (m DotnetManager) ManagerVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "dotnet", "--version").Output()
	if err != nil {
		return "", err
	}
	return "v" + strings.TrimSpace(string(out)), nil
}
