package golang

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"senk-tui/managers"
	"sort"
	"strings"
	"time"

	"fmt"
)

type GoManager struct{}

func (m GoManager) Name() string { return "Go (GOPATH)" }

func (m GoManager) IsInstalled() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

func (m GoManager) Fetch() ([]managers.Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	gopathCmd := exec.CommandContext(ctx, "go", "env", "GOPATH")
	gopathBytes, err := gopathCmd.Output()
	if err != nil {
		return nil, err
	}
	gopath := strings.TrimSpace(string(gopathBytes))
	if gopath == "" {
		home, _ := os.UserHomeDir()
		gopath = filepath.Join(home, "go")
	}

	binDir := filepath.Join(gopath, "bin")
	entries, err := os.ReadDir(binDir)
	if err != nil {
		return []managers.Dependency{}, nil
	}

	var binaries []string
	for _, e := range entries {
		if !e.IsDir() {
			binaries = append(binaries, filepath.Join(binDir, e.Name()))
		}
	}

	if len(binaries) == 0 {
		return []managers.Dependency{}, nil
	}

	args := append([]string{"version", "-m"}, binaries...)
	versionCmd := exec.CommandContext(ctx, "go", args...)
	out, _ := versionCmd.Output()

	return parseGoBinOutput(out)
}

func parseGoBinOutput(data []byte) ([]managers.Dependency, error) {
	var deps []managers.Dependency
	lines := strings.Split(string(data), "\n")

	var currentName string
	var currentVersion string

	for _, line := range lines {
		if !strings.HasPrefix(line, "\t") {
			parts := strings.Split(line, ":")
			if len(parts) > 0 {
				pathParts := strings.Split(parts[0], "/")
				if len(pathParts) > 0 {
					name := pathParts[len(pathParts)-1]
					currentName = strings.TrimSpace(name)
					currentVersion = "(devel)"
				}
			}
		} else if strings.HasPrefix(line, "\tmod\t") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				currentVersion = parts[2]
			}
			if currentName != "" {
				deps = append(deps, managers.Dependency{
					Name:    currentName,
					Version: currentVersion,
				})
				currentName = ""
			}
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

func (m GoManager) ManagerVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "go", "env", "GOVERSION").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Helper to get file size
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

func (m GoManager) FetchDetails(dep managers.Dependency) (managers.PackageDetails, error) {
	details := managers.PackageDetails{Name: dep.Name, Version: dep.Version}

	// Re-construct the path
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, "go", "bin", dep.Name)
	details.Path = path
	details.Size = getFileSizeStr(path)

	return details, nil
}
func (m GoManager) Uninstall(dep managers.Dependency) error {
	home, _ := os.UserHomeDir()
	return os.Remove(filepath.Join(home, "go", "bin", dep.Name))
}
