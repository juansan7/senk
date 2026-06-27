package rust

import (
	"os/exec"
	"senk-tui/managers"
	"strings"

	"context"
	"time"
)

type RustLanguage struct{}

func (l RustLanguage) Name() string { return "Rust" }
func (l RustLanguage) IsInstalled() bool {
	_, err := exec.LookPath("rustc")
	return err == nil
}
func (l RustLanguage) LanguageVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "rustc", "--version").Output()
	if err != nil {
		return "", err
	}
	parts := strings.Fields(string(out))
	if len(parts) >= 2 {
		return "v" + parts[1], nil
	}
	return "", nil
}
func (l RustLanguage) Managers() []managers.PackageManager {
	return []managers.PackageManager{CargoManager{}}
}
