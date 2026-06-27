package managers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

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

// ----------------- Go -----------------
func (m GoManager) FetchDetails(dep Dependency) (PackageDetails, error) {
	details := PackageDetails{Name: dep.Name, Version: dep.Version}

	// Re-construct the path
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, "go", "bin", dep.Name)
	details.Path = path
	details.Size = getFileSizeStr(path)

	return details, nil
}
func (m GoManager) Uninstall(dep Dependency) error {
	home, _ := os.UserHomeDir()
	return os.Remove(filepath.Join(home, "go", "bin", dep.Name))
}

// ----------------- Cargo -----------------
func (m CargoManager) FetchDetails(dep Dependency) (PackageDetails, error) {
	details := PackageDetails{Name: dep.Name, Version: dep.Version}
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".cargo", "bin", dep.Name)
	details.Path = path
	details.Size = getFileSizeStr(path)
	return details, nil
}
func (m CargoManager) Uninstall(dep Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "cargo", "uninstall", dep.Name).Run()
}

// ----------------- Gem -----------------
func (m GemManager) FetchDetails(dep Dependency) (PackageDetails, error) {
	return PackageDetails{Name: dep.Name, Version: dep.Version}, nil
}
func (m GemManager) Uninstall(dep Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "gem", "uninstall", "-a", "-x", dep.Name).Run()
}

// ----------------- Dotnet -----------------
func (m DotnetManager) FetchDetails(dep Dependency) (PackageDetails, error) {
	return PackageDetails{Name: dep.Name, Version: dep.Version}, nil
}
func (m DotnetManager) Uninstall(dep Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "dotnet", "tool", "uninstall", "-g", dep.Name).Run()
}

// ----------------- Composer -----------------
func (m ComposerManager) FetchDetails(dep Dependency) (PackageDetails, error) {
	return PackageDetails{Name: dep.Name, Version: dep.Version}, nil
}
func (m ComposerManager) Uninstall(dep Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "composer", "global", "remove", dep.Name).Run()
}

// ----------------- Bun -----------------
func (m BunManager) FetchDetails(dep Dependency) (PackageDetails, error) {
	return PackageDetails{Name: dep.Name, Version: dep.Version}, nil
}
func (m BunManager) Uninstall(dep Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "bun", "pm", "rm", "-g", dep.Name).Run()
}

// ----------------- Pub -----------------
func (m PubManager) FetchDetails(dep Dependency) (PackageDetails, error) {
	return PackageDetails{Name: dep.Name, Version: dep.Version}, nil
}
func (m PubManager) Uninstall(dep Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "dart", "pub", "global", "deactivate", dep.Name).Run()
}

// ----------------- Luarocks -----------------
func (m LuarocksManager) FetchDetails(dep Dependency) (PackageDetails, error) {
	return PackageDetails{Name: dep.Name, Version: dep.Version}, nil
}
func (m LuarocksManager) Uninstall(dep Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "luarocks", "remove", "--force", dep.Name).Run()
}
