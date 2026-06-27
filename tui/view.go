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

	paneWidth1 := (m.Width*25)/100 - 2
	paneWidth2 := (m.Width*25)/100 - 2
	paneWidth3 := m.Width - paneWidth1 - paneWidth2 - 6
	paneHeight := m.Height - 4

	leftStyle := baseStyle
	midStyle := baseStyle
	rightStyle := baseStyle

	if m.Focus == FocusLeft {
		leftStyle = activeStyle
	} else if m.Focus == FocusMiddle {
		midStyle = activeStyle
	} else {
		rightStyle = activeStyle
	}

	leftStyle = leftStyle.Width(paneWidth1).Height(paneHeight)
	midStyle = midStyle.Width(paneWidth2).Height(paneHeight)
	rightStyle = rightStyle.Width(paneWidth3).Height(paneHeight)

	leftContent := m.viewLeftPane()
	midContent := m.viewMiddlePane()
	rightContent := m.viewRightPane(paneWidth3, paneHeight-2)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(leftContent),
		midStyle.Render(midContent),
		rightStyle.Render(rightContent),
	)
}

func (m Model) viewLeftPane() string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("Languages") + "\n\n")

	for i, lang := range m.Languages {
		cursor := "  "
		style := unselectedItemStyle
		if i == m.LangIndex {
			cursor = "> "
			style = selectedItemStyle
		}

		name := lang.Language.Name()
		if lang.Version != "" {
			name = fmt.Sprintf("%s (%s)", name, lang.Version)
		}
		s.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(name)))
	}
	return s.String()
}

func (m Model) viewMiddlePane() string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("Managers") + "\n\n")

	if len(m.Languages) == 0 {
		return s.String()
	}

	lang := m.Languages[m.LangIndex]
	for i, mgrData := range lang.Managers {
		cursor := "  "
		style := unselectedItemStyle
		if i == m.MgrIndex {
			cursor = "> "
			style = selectedItemStyle
		}

		status := ""
		switch mgrData.State {
		case StateLoading:
			status = m.Spinner.View()
		case StateDone:
			status = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓")
		case StateError:
			status = errorStyle.Render("✗")
		}

		name := mgrData.Manager.Name()
		if mgrData.Version != "" {
			name = fmt.Sprintf("%s (%s)", name, mgrData.Version)
		}

		s.WriteString(fmt.Sprintf("%s%s [%s]\n", cursor, style.Render(name), status))
	}
	return s.String()
}

func (m Model) viewRightPane(width int, maxVisible int) string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("Packages") + "\n\n")

	if len(m.Languages) == 0 || len(m.Languages[m.LangIndex].Managers) == 0 {
		s.WriteString("No packages available.")
		return s.String()
	}

	activeManager := m.Languages[m.LangIndex].Managers[m.MgrIndex]

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

	startIdx := 0
	endIdx := len(activeManager.Dependencies)

	if m.DepIndex >= maxVisible {
		startIdx = m.DepIndex - maxVisible + 1
	}
	if endIdx > startIdx+maxVisible {
		endIdx = startIdx + maxVisible
	}

	for i := startIdx; i < endIdx; i++ {
		dep := activeManager.Dependencies[i]
		cursor := "  "
		style := unselectedItemStyle

		if m.Focus == FocusRight && i == m.DepIndex {
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
