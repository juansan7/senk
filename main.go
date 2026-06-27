package main

import (
	"fmt"
	"os"

	"senk-tui/managers"
	"senk-tui/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	allLanguages := []managers.Language{
		managers.MacLanguage{},
		managers.NodeLanguage{},
		managers.PythonLanguage{},
		managers.GolangLanguage{},
		managers.RustLanguage{},
		managers.RubyLanguage{},
		managers.DotnetLanguage{},
		managers.PhpLanguage{},
		managers.BunLanguage{},
		managers.DartLanguage{},
		managers.LuaLanguage{},
	}

	var installedLanguages []managers.Language
	for _, l := range allLanguages {
		if l.IsInstalled() {
			installedLanguages = append(installedLanguages, l)
		}
	}

	model := tui.NewModel(installedLanguages)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running senk-tui: %v", err)
		os.Exit(1)
	}
}
