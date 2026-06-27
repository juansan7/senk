package main

import (
	"fmt"
	"os"

	"senk-tui/managers"
	"senk-tui/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var availableManagers []managers.PackageManager
	
	allManagers := []managers.PackageManager{
		managers.BrewManager{},
		managers.NpmManager{},
		managers.PipManager{},
		managers.GoManager{},
		managers.CargoManager{},
		managers.GemManager{},
		managers.DotnetManager{},
		managers.ComposerManager{},
		managers.BunManager{},
		managers.PubManager{},
		managers.LuarocksManager{},
	}
	
	for _, m := range allManagers {
		if m.IsInstalled() {
			availableManagers = append(availableManagers, m)
		}
	}
	
	model := tui.NewModel(availableManagers)
	
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running senk-tui: %v", err)
		os.Exit(1)
	}
}
