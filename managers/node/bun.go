package node

import (
	"context"
	"os/exec"
	"regexp"
	"senk-tui/managers"
	"sort"
	"strings"
	"time"
)

type BunManager struct{}

func (m BunManager) Name() string { return "Bun" }
func (m BunManager) IsInstalled() bool {
	_, err := exec.LookPath("bun")
	return err == nil
}
func (m BunManager) Fetch() ([]managers.Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "bun", "pm", "ls", "-g").Output()
	if len(out) > 0 {
		return parseBunOutput(out)
	}
	return nil, err
}

func parseBunOutput(data []byte) ([]managers.Dependency, error) {
	var deps []managers.Dependency
	lines := strings.Split(string(data), "\n")
	re := regexp.MustCompile(`(?:├──|└──)\s+(.+)@(.+)`)
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			deps = append(deps, managers.Dependency{Name: matches[1], Version: matches[2]})
		}
	}
	sort.Slice(deps, func(i, j int) bool { return deps[i].Name < deps[j].Name })
	if deps == nil {
		deps = []managers.Dependency{}
	}
	return deps, nil
}

func (m BunManager) ManagerVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "bun", "--version").Output()
	if err != nil {
		return "", err
	}
	return "v" + strings.TrimSpace(string(out)), nil
}

func (m BunManager) FetchDetails(dep managers.Dependency) (managers.PackageDetails, error) {
	return managers.PackageDetails{Name: dep.Name, Version: dep.Version}, nil
}
func (m BunManager) Uninstall(dep managers.Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "bun", "pm", "rm", "-g", dep.Name).Run()
}
