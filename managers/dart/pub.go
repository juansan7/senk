package dart

import (
	"context"
	"os/exec"
	"senk-tui/managers"
	"sort"
	"strings"
	"time"
)

type PubManager struct{}

func (m PubManager) Name() string { return "Dart (pub)" }
func (m PubManager) IsInstalled() bool {
	_, err := exec.LookPath("dart")
	return err == nil
}
func (m PubManager) Fetch() ([]managers.Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "dart", "pub", "global", "list").Output()
	if len(out) > 0 {
		return parsePubOutput(out)
	}
	return nil, err
}

func parsePubOutput(data []byte) ([]managers.Dependency, error) {
	var deps []managers.Dependency
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			deps = append(deps, managers.Dependency{Name: parts[0], Version: parts[1]})
		}
	}
	sort.Slice(deps, func(i, j int) bool { return deps[i].Name < deps[j].Name })
	if deps == nil {
		deps = []managers.Dependency{}
	}
	return deps, nil
}

func (m PubManager) ManagerVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "dart", "--version").Output()
	if err != nil {
		return "", err
	}
	parts := strings.Fields(string(out))
	if len(parts) >= 4 {
		return "v" + parts[3], nil
	}
	return "", nil
}

func (m PubManager) FetchDetails(dep managers.Dependency) (managers.PackageDetails, error) {
	return managers.PackageDetails{Name: dep.Name, Version: dep.Version}, nil
}
func (m PubManager) Uninstall(dep managers.Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "dart", "pub", "global", "deactivate", dep.Name).Run()
}
