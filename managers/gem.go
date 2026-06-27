package managers

import (
	"context"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
)

type GemManager struct{}

func (m GemManager) Name() string { return "Ruby (gem)" }
func (m GemManager) IsInstalled() bool {
	_, err := exec.LookPath("gem")
	return err == nil
}
func (m GemManager) Fetch() ([]Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "gem", "list", "--local").Output()
	if len(out) > 0 {
		return parseGemOutput(out)
	}
	return nil, err
}

func parseGemOutput(data []byte) ([]Dependency, error) {
	var deps []Dependency
	lines := strings.Split(string(data), "\n")
	re := regexp.MustCompile(`^([^ ]+)\s+\(([^)]+)\)`)
	for _, line := range lines {
		matches := re.FindStringSubmatch(strings.TrimSpace(line))
		if len(matches) == 3 {
			versionsRaw := matches[2]
			var userVersions []string
			for _, v := range strings.Split(versionsRaw, ",") {
				v = strings.TrimSpace(v)
				if !strings.HasPrefix(v, "default:") {
					userVersions = append(userVersions, v)
				}
			}
			if len(userVersions) > 0 {
				deps = append(deps, Dependency{Name: matches[1], Version: strings.Join(userVersions, ", ")})
			}
		}
	}
	sort.Slice(deps, func(i, j int) bool { return deps[i].Name < deps[j].Name })
	if deps == nil {
		deps = []Dependency{}
	}
	return deps, nil
}

func (m GemManager) ManagerVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "gem", "--version").Output()
	if err != nil {
		return "", err
	}
	return "v" + strings.TrimSpace(string(out)), nil
}
