package dotnet

import (
	"os/exec"
	"senk-tui/managers"
	"strings"

	"context"
	"time"
)

type DotnetLanguage struct{}

func (l DotnetLanguage) Name() string { return ".NET" }
func (l DotnetLanguage) IsInstalled() bool {
	_, err := exec.LookPath("dotnet")
	return err == nil
}
func (l DotnetLanguage) LanguageVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "dotnet", "--version").Output()
	if err != nil {
		return "", err
	}
	return "v" + strings.TrimSpace(string(out)), nil
}
func (l DotnetLanguage) Managers() []managers.PackageManager {
	return []managers.PackageManager{DotnetManager{}}
}
