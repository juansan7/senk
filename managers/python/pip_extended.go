package python

import (
	"context"
	"os/exec"
	"senk-tui/managers"
	"strings"
	"time"
)

func (m PipManager) FetchDetails(dep managers.Dependency) (managers.PackageDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	details := managers.PackageDetails{Name: dep.Name, Version: dep.Version}

	out, err := exec.CommandContext(ctx, m.getExecutable(), "show", dep.Name).Output()
	if err != nil {
		return details, nil
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Summary: ") {
			details.Description = strings.TrimPrefix(line, "Summary: ")
		} else if strings.HasPrefix(line, "Author: ") {
			details.Author = strings.TrimPrefix(line, "Author: ")
		} else if strings.HasPrefix(line, "Location: ") {
			details.Path = strings.TrimPrefix(line, "Location: ")
		}
	}
	return details, nil
}

func (m PipManager) Uninstall(dep managers.Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// -y bypasses the "Are you sure?" prompt
	return exec.CommandContext(ctx, m.getExecutable(), "uninstall", "-y", dep.Name).Run()
}
