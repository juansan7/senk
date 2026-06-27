package lua

import (
	"os/exec"
	"senk-tui/managers"
	"strings"

	"context"
	"time"
)

type LuaLanguage struct{}

func (l LuaLanguage) Name() string { return "Lua" }
func (l LuaLanguage) IsInstalled() bool {
	_, err := exec.LookPath("lua")
	return err == nil
}
func (l LuaLanguage) LanguageVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "lua", "-v").Output()
	if err != nil {
		return "", err
	}
	parts := strings.Fields(string(out))
	if len(parts) >= 2 {
		return "v" + parts[1], nil
	}
	return "", nil
}
func (l LuaLanguage) Managers() []managers.PackageManager {
	return []managers.PackageManager{LuarocksManager{}}
}
