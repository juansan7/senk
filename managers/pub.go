package managers

import (
	"context"
	"os/exec"
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
func (m PubManager) Fetch() ([]Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "dart", "pub", "global", "list").Output()
	if len(out) > 0 {
		return parsePubOutput(out)
	}
	return nil, err
}

func parsePubOutput(data []byte) ([]Dependency, error) {
	var deps []Dependency
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
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
