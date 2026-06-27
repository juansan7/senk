package node

import (
	"os/exec"
	"senk-tui/managers"
	"strings"

	"context"
	"time"
)

type NodeLanguage struct{}

func (l NodeLanguage) Name() string { return "Node.js" }
func (l NodeLanguage) IsInstalled() bool {
	_, err := exec.LookPath("node")
	return err == nil
}
func (l NodeLanguage) LanguageVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "node", "--version").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
func (l NodeLanguage) Managers() []managers.PackageManager {
	return []managers.PackageManager{NpmManager{}} // We can add YarnManager{}, PnpmManager{} later
}

type BunLanguage struct{}

func (l BunLanguage) Name() string { return "Bun" }
func (l BunLanguage) IsInstalled() bool {
	_, err := exec.LookPath("bun")
	return err == nil
}
func (l BunLanguage) LanguageVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "bun", "--version").Output()
	if err != nil {
		return "", err
	}
	return "v" + strings.TrimSpace(string(out)), nil
}
func (l BunLanguage) Managers() []managers.PackageManager {
	return []managers.PackageManager{BunManager{}}
}
