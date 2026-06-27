package managers

import (
	"context"
	"encoding/json"
	"os/exec"
	"time"
)

// BrewInfoJSON represents the output of `brew info <pkg> --json`
type BrewInfoJSON []struct {
	Desc     string `json:"desc"`
	Homepage string `json:"homepage"`
	Tap      string `json:"tap"`
	KegOnly  bool   `json:"keg_only"`
}

func (m BrewManager) FetchDetails(dep Dependency) (PackageDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	details := PackageDetails{
		Name:    dep.Name,
		Version: dep.Version,
	}

	// For brew, we use `brew info <pkg> --json`
	out, err := exec.CommandContext(ctx, "brew", "info", dep.Name, "--json").Output()
	if err != nil {
		return details, nil // Return basic info if brew info fails
	}

	var info BrewInfoJSON
	if err := json.Unmarshal(out, &info); err == nil && len(info) > 0 {
		details.Description = info[0].Desc
		details.Homepage = info[0].Homepage
		details.Author = info[0].Tap
	}

	return details, nil
}

func (m BrewManager) Uninstall(dep Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Using --force to bypass any safety prompts, making it TUI friendly
	cmd := exec.CommandContext(ctx, "brew", "uninstall", "--force", dep.Name)
	return cmd.Run()
}
