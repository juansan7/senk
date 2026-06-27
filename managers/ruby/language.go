package ruby

import (
	"os/exec"
	"senk-tui/managers"
	"strings"

	"context"
	"time"
)

type RubyLanguage struct{}

func (l RubyLanguage) Name() string { return "Ruby" }
func (l RubyLanguage) IsInstalled() bool {
	_, err := exec.LookPath("ruby")
	return err == nil
}
func (l RubyLanguage) LanguageVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "ruby", "--version").Output()
	if err != nil {
		return "", err
	}
	parts := strings.Fields(string(out))
	if len(parts) >= 2 {
		return "v" + parts[1], nil
	}
	return "", nil
}
func (l RubyLanguage) Managers() []managers.PackageManager {
	return []managers.PackageManager{GemManager{}}
}
