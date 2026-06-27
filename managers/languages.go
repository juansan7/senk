package managers

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// ---------------------------------------------------------
// Node.js
// ---------------------------------------------------------
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
func (l NodeLanguage) Managers() []PackageManager {
	return []PackageManager{NpmManager{}} // We can add YarnManager{}, PnpmManager{} later
}

// ---------------------------------------------------------
// Python
// ---------------------------------------------------------
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
func (l PythonLanguage) Managers() []PackageManager {
	return []PackageManager{PipManager{}}
}

// ---------------------------------------------------------
// Rust
// ---------------------------------------------------------
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
func (l RustLanguage) Managers() []PackageManager {
	return []PackageManager{CargoManager{}}
}

// ---------------------------------------------------------
// Go
// ---------------------------------------------------------
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
func (l GolangLanguage) Managers() []PackageManager {
	return []PackageManager{GoManager{}}
}

// ---------------------------------------------------------
// Ruby
// ---------------------------------------------------------
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
func (l RubyLanguage) Managers() []PackageManager {
	return []PackageManager{GemManager{}}
}

// ---------------------------------------------------------
// .NET
// ---------------------------------------------------------
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
func (l DotnetLanguage) Managers() []PackageManager {
	return []PackageManager{DotnetManager{}}
}

// ---------------------------------------------------------
// PHP
// ---------------------------------------------------------
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
func (l PhpLanguage) Managers() []PackageManager {
	return []PackageManager{ComposerManager{}}
}

// ---------------------------------------------------------
// Dart
// ---------------------------------------------------------
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
func (l DartLanguage) Managers() []PackageManager {
	return []PackageManager{PubManager{}}
}

// ---------------------------------------------------------
// Lua
// ---------------------------------------------------------
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
func (l LuaLanguage) Managers() []PackageManager {
	return []PackageManager{LuarocksManager{}}
}

// ---------------------------------------------------------
// Bun (Runtime & Manager)
// ---------------------------------------------------------
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
func (l BunLanguage) Managers() []PackageManager {
	return []PackageManager{BunManager{}}
}

// ---------------------------------------------------------
// macOS (System)
// ---------------------------------------------------------
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
func (l MacLanguage) Managers() []PackageManager {
	return []PackageManager{BrewManager{}}
}
