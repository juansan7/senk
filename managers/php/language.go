package php

import (
	"os/exec"
	"senk-tui/managers"
	"strings"

	"context"
	"time"
)

type PhpLanguage struct{}

func (l PhpLanguage) Name() string { return "PHP" }
func (l PhpLanguage) IsInstalled() bool {
	_, err := exec.LookPath("php")
	return err == nil
}
func (l PhpLanguage) LanguageVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "php", "--version").Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 {
			return "v" + parts[1], nil
		}
	}
	return "", nil
}
func (l PhpLanguage) Managers() []managers.PackageManager {
	return []managers.PackageManager{ComposerManager{}}
}
