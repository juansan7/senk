package tui

import "github.com/charmbracelet/lipgloss"

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	activeStyle = baseStyle.
			BorderForeground(lipgloss.Color("62")) // Purpleish-blue border when focused

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("62")).
			PaddingLeft(1).
			PaddingRight(1)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("212")). // Pink text when selected
				Bold(true)

	unselectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	versionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")) // Dark grey for versions

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			Background(lipgloss.Color("235"))

	dangerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
)
