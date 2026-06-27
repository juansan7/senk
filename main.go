package main

import (
	"fmt"
	"os"

	"senk-tui/managers"
	"senk-tui/managers/dart"
	"senk-tui/managers/dotnet"
	"senk-tui/managers/golang"
	"senk-tui/managers/lua"
	"senk-tui/managers/mac"
	"senk-tui/managers/node"
	"senk-tui/managers/php"
	"senk-tui/managers/python"
	"senk-tui/managers/ruby"
	"senk-tui/managers/rust"
	"senk-tui/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	allLanguages := []managers.Language{
		mac.MacLanguage{},
		node.NodeLanguage{},
		python.PythonLanguage{},
		golang.GolangLanguage{},
		rust.RustLanguage{},
		ruby.RubyLanguage{},
		dotnet.DotnetLanguage{},
		php.PhpLanguage{},
		node.BunLanguage{},
		dart.DartLanguage{},
		lua.LuaLanguage{},
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
