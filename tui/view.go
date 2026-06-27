package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Initializing senk-tui..."
	}

	// Layout constraints
	leftPaneWidth := (m.Width / 3) - 2
	rightPaneWidth := m.Width - leftPaneWidth - 4
	paneHeight := m.Height - 4 // Account for borders

	leftStyle := baseStyle
	rightStyle := baseStyle

	if m.Focus == FocusLeft {
		leftStyle = activeStyle
	} else {
		rightStyle = activeStyle
	}

	leftStyle = leftStyle.Width(leftPaneWidth).Height(paneHeight)
	rightStyle = rightStyle.Width(rightPaneWidth).Height(paneHeight)

	leftContent := m.viewLeftPane()
	rightContent := m.viewRightPane(rightPaneWidth, paneHeight-2)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(leftContent),
		rightStyle.Render(rightContent),
	)
}

func (m Model) viewLeftPane() string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("Managers") + "\n\n")

	for i, manager := range m.Managers {
		cursor := "  "
		style := unselectedItemStyle
		if i == m.LeftIndex {
			cursor = "> "
			style = selectedItemStyle
		}

		status := ""
		switch manager.State {
		case StateLoading:
			status = m.Spinner.View()
		case StateDone:
			status = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓")
		case StateError:
			status = errorStyle.Render("✗")
		}

		name := manager.Manager.Name()
		if manager.Version != "" {
			name = fmt.Sprintf("%s (%s)", name, manager.Version)
		}

		s.WriteString(fmt.Sprintf("%s%s [%s]\n", cursor, style.Render(name), status))
	}
	return s.String()
}

func (m Model) viewRightPane(width int, maxVisible int) string {
	if len(m.Managers) == 0 {
		return "No package managers found on system."
	}

	activeManager := m.Managers[m.LeftIndex]

	var s strings.Builder
	s.WriteString(titleStyle.Render(activeManager.Manager.Name()+" Packages") + "\n\n")

	switch activeManager.State {
	case StateLoading:
		s.WriteString(fmt.Sprintf("%s Fetching dependencies...\nThis might take a moment.", m.Spinner.View()))
		return s.String()
	case StateError:
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", activeManager.Err)))
		return s.String()
	}

	if len(activeManager.Dependencies) == 0 {
		s.WriteString("No globally installed packages found.")
		return s.String()
	}

	// Simple pagination/scrolling window
	startIdx := 0
	endIdx := len(activeManager.Dependencies)
	
	if m.RightIndex >= maxVisible {
		startIdx = m.RightIndex - maxVisible + 1
	}
	if endIdx > startIdx+maxVisible {
		endIdx = startIdx + maxVisible
	}

	for i := startIdx; i < endIdx; i++ {
		dep := activeManager.Dependencies[i]
		cursor := "  "
		style := unselectedItemStyle

		if m.Focus == FocusRight && i == m.RightIndex {
			cursor = "> "
			style = selectedItemStyle
		}

		nameRender := style.Render(dep.Name)
		versionRender := versionStyle.Render(dep.Version)
		
		nameLen := lipgloss.Width(nameRender)
		verLen := lipgloss.Width(versionRender)
		
		padding := width - nameLen - verLen - 6
		if padding < 1 {
			padding = 1
		}
		padStr := strings.Repeat(".", padding)
		padRender := lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(padStr)

		s.WriteString(fmt.Sprintf("%s%s %s %s\n", cursor, nameRender, padRender, versionRender))
	}

	return s.String()
}
