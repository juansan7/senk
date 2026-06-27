package python

import (
	"os/exec"
	"senk-tui/managers"
	"strings"

	"context"
	"time"
)

type PythonLanguage struct{}

func (l PythonLanguage) Name() string { return "Python" }
func (l PythonLanguage) IsInstalled() bool {
	if _, err := exec.LookPath("python3"); err == nil {
		return true
	}
	_, err := exec.LookPath("python")
	return err == nil
}
func (l PythonLanguage) LanguageVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bin := "python"
	if _, err := exec.LookPath("python3"); err == nil {
		bin = "python3"
	}
	out, err := exec.CommandContext(ctx, bin, "--version").Output()
	if err != nil {
		return "", err
	}
	parts := strings.Fields(string(out))
	if len(parts) >= 2 {
		return "v" + parts[1], nil
	}
	return "", nil
}
func (l PythonLanguage) Managers() []managers.PackageManager {
	return []managers.PackageManager{PipManager{}}
}
