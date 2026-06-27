package managers

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type GoManager struct{}

func (m GoManager) Name() string { return "Go (GOPATH)" }

func (m GoManager) IsInstalled() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

func (m GoManager) Fetch() ([]Dependency, error) {
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
		return []Dependency{}, nil
	}

	var binaries []string
	for _, e := range entries {
		if !e.IsDir() {
			binaries = append(binaries, filepath.Join(binDir, e.Name()))
		}
	}

	if len(binaries) == 0 {
		return []Dependency{}, nil
	}

	args := append([]string{"version", "-m"}, binaries...)
	versionCmd := exec.CommandContext(ctx, "go", args...)
	out, _ := versionCmd.Output()

	return parseGoBinOutput(out)
}

func parseGoBinOutput(data []byte) ([]Dependency, error) {
	var deps []Dependency
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
				deps = append(deps, Dependency{
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
		deps = []Dependency{}
	}

	return deps, nil
}
