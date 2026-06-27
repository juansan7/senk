package golang

import (
	"os/exec"
	"senk-tui/managers"
	"strings"

	"context"
	"time"
)

type GolangLanguage struct{}

func (l GolangLanguage) Name() string { return "Go" }
func (l GolangLanguage) IsInstalled() bool {
	_, err := exec.LookPath("go")
	return err == nil
}
func (l GolangLanguage) LanguageVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "go", "env", "GOVERSION").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
func (l GolangLanguage) Managers() []managers.PackageManager {
	return []managers.PackageManager{GoManager{}}
}
