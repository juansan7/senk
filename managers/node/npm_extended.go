package node

import (
	"context"
	"encoding/json"
	"os/exec"
	"senk-tui/managers"
	"time"
)

type NpmViewJSON struct {
	Description string      `json:"description"`
	Homepage    string      `json:"homepage"`
	Author      interface{} `json:"author"` // can be string or object
}

func (m NpmManager) FetchDetails(dep managers.Dependency) (managers.PackageDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	details := managers.PackageDetails{Name: dep.Name, Version: dep.Version}

	out, err := exec.CommandContext(ctx, "npm", "view", dep.Name, "--json").Output()
	if err != nil {
		return details, nil
	}

	var view NpmViewJSON
	if err := json.Unmarshal(out, &view); err == nil {
		details.Description = view.Description
		details.Homepage = view.Homepage
		switch a := view.Author.(type) {
		case string:
			details.Author = a
		case map[string]interface{}:
			if name, ok := a["name"].(string); ok {
				details.Author = name
			}
		}
	}
	return details, nil
}

func (m NpmManager) Uninstall(dep managers.Dependency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "npm", "uninstall", "-g", dep.Name).Run()
}
