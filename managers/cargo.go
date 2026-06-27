package managers

import (
	"context"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
)

type CargoManager struct{}

func (m CargoManager) Name() string { return "Rust (Cargo)" }

func (m CargoManager) IsInstalled() bool {
	_, err := exec.LookPath("cargo")
	return err == nil
}

func (m CargoManager) Fetch() ([]Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "cargo", "install", "--list")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseCargoOutput(out)
}

func parseCargoOutput(data []byte) ([]Dependency, error) {
	var deps []Dependency
	lines := strings.Split(string(data), "\n")

	re := regexp.MustCompile(`^([^ ]+) v([^:]+):`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			deps = append(deps, Dependency{
				Name:    matches[1],
				Version: "v" + matches[2],
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
