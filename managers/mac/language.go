package mac

import (
	"os/exec"
	"senk-tui/managers"
	"strings"

	"context"
	"time"
)

type MacLanguage struct{}

func (l MacLanguage) Name() string { return "macOS" }
func (l MacLanguage) IsInstalled() bool {
	_, err := exec.LookPath("brew") // Brew implies macOS/Linux system management
	return err == nil
}
func (l MacLanguage) LanguageVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "sw_vers", "-productVersion").Output()
	if err != nil {
		return "", nil
	} // ignore error on linux
	return "v" + strings.TrimSpace(string(out)), nil
}
func (l MacLanguage) Managers() []managers.PackageManager {
	return []managers.PackageManager{BrewManager{}}
}
