package php

import (
	"context"
	"encoding/json"
	"os/exec"
	"senk-tui/managers"
	"sort"
	"strings"
	"time"
)

type ComposerManager struct{}

func (m ComposerManager) Name() string { return "PHP (composer)" }
func (m ComposerManager) IsInstalled() bool {
	_, err := exec.LookPath("composer")
	return err == nil
}
func (m ComposerManager) Fetch() ([]managers.Dependency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "composer", "global", "show", "--direct", "--format=json").Output()
	if len(out) > 0 {
		return parseComposerOutput(out)
	}
	return nil, err
}

type composerOutput struct {
	Installed []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"installed"`
}

func parseComposerOutput(data []byte) ([]managers.Dependency, error) {
	var out composerOutput
	if len(data) > 0 {
		_ = json.Unmarshal(data, &out)
	}
	var deps []managers.Dependency
	for _, pkg := range out.Installed {
		deps = append(deps, managers.Dependency{Name: pkg.Name, Version: pkg.Version})
	}
	sort.Slice(deps, func(i, j int) bool { return deps[i].Name < deps[j].Name })
	if deps == nil {
		deps = []managers.Dependency{}
	}
	return deps, nil
}

func (m ComposerManager) ManagerVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "composer", "--version").Output()
	if err != nil {
		return "", err
	}
	parts := strings.Fields(string(out))
	if len(parts) >= 3 {
		return "v" + parts[2], nil
	}
	return "", nil
}

func (m ComposerManager) FetchDetails(dep managers.Dependency) (managers.PackageDetails, error) {
	return managers.PackageDetails{Name: dep.Name, Version: dep.Version}, nil
}
func (m ComposerManager) Uninstall(dep managers.Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "composer", "global", "remove", dep.Name).Run()
}
