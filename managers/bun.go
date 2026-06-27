package managers

import (
	"context"
	"os/exec"
	"regexp"
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
func (m BunManager) Fetch() ([]Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "bun", "pm", "ls", "-g").Output()
	if len(out) > 0 { return parseBunOutput(out) }
	return nil, err
}

func parseBunOutput(data []byte) ([]Dependency, error) {
	var deps []Dependency
	lines := strings.Split(string(data), "\n")
	re := regexp.MustCompile(`(?:├──|└──)\s+(.+)@(.+)`)
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			deps = append(deps, Dependency{Name: matches[1], Version: matches[2]})
		}
	}
	sort.Slice(deps, func(i, j int) bool { return deps[i].Name < deps[j].Name })
	if deps == nil { deps = []Dependency{} }
	return deps, nil
}
