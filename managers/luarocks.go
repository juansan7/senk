package managers

import (
	"context"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type LuarocksManager struct{}

func (m LuarocksManager) Name() string { return "Lua (luarocks)" }
func (m LuarocksManager) IsInstalled() bool {
	_, err := exec.LookPath("luarocks")
	return err == nil
}
func (m LuarocksManager) Fetch() ([]Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "luarocks", "list", "--porcelain").Output()
	if len(out) > 0 {
		return parseLuarocksOutput(out)
	}
	return nil, err
}

func parseLuarocksOutput(data []byte) ([]Dependency, error) {
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

func (m LuarocksManager) ManagerVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "luarocks", "--version").Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 {
			return "v" + parts[1], nil
		}
	}
	return "", nil
}
