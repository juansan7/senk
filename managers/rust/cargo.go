package rust

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"senk-tui/managers"
	"sort"
	"strings"
	"time"

	"os"
	"path/filepath"
)

type CargoManager struct{}

func (m CargoManager) Name() string { return "Rust (Cargo)" }

func (m CargoManager) IsInstalled() bool {
	_, err := exec.LookPath("cargo")
	return err == nil
}

func (m CargoManager) Fetch() ([]managers.Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "cargo", "install", "--list")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseCargoOutput(out)
}

func parseCargoOutput(data []byte) ([]managers.Dependency, error) {
	var deps []managers.Dependency
	lines := strings.Split(string(data), "\n")

	re := regexp.MustCompile(`^([^ ]+) v([^:]+):`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			deps = append(deps, managers.Dependency{
				Name:    matches[1],
				Version: "v" + matches[2],
			})
		}
	}

	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Name < deps[j].Name
	})

	if deps == nil {
		deps = []managers.Dependency{}
	}

	return deps, nil
}

func (m CargoManager) ManagerVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "cargo", "--version").Output()
	if err != nil {
		return "", err
	}
	parts := strings.Fields(string(out))
	if len(parts) >= 2 {
		return "v" + parts[1], nil
	}
	return "", nil
}

func (m CargoManager) FetchDetails(dep managers.Dependency) (managers.PackageDetails, error) {
	details := managers.PackageDetails{Name: dep.Name, Version: dep.Version}
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".cargo", "bin", dep.Name)
	details.Path = path
	details.Size = getFileSizeStr(path)
	return details, nil
}
func (m CargoManager) Uninstall(dep managers.Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "cargo", "uninstall", dep.Name).Run()
}

func getFileSizeStr(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}
	size := info.Size()
	if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024.0)
	}
	return fmt.Sprintf("%.2f MB", float64(size)/(1024.0*1024.0))
}
