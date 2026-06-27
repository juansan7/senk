package dart

import (
	"os/exec"
	"senk-tui/managers"
	"strings"

	"context"
	"time"
)

type DartLanguage struct{}

func (l DartLanguage) Name() string { return "Dart" }
func (l DartLanguage) IsInstalled() bool {
	_, err := exec.LookPath("dart")
	return err == nil
}
func (l DartLanguage) LanguageVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "dart", "--version").Output()
	if err != nil {
		return "", err
	}
	parts := strings.Fields(string(out))
	if len(parts) >= 4 {
		return "v" + parts[3], nil
	}
	return "", nil
}
func (l DartLanguage) Managers() []managers.PackageManager {
	return []managers.PackageManager{PubManager{}}
}
