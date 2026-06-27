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

	leftContent := m.viewLeftPane(paneWidth1)
	midContent := m.viewMiddlePane(paneWidth2)
	rightContent := m.viewRightPane(paneWidth3, paneHeight-2)

	ui := lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(leftContent),
		midStyle.Render(midContent),
		rightStyle.Render(rightContent),
	)

	if m.Modal.State != ModalClosed {
		return m.renderModalOverlay(ui)
	}

	return ui
}

func (m Model) renderModalOverlay(bg string) string {
	var s strings.Builder

	switch m.Modal.State {
	case ModalLoading:
		s.WriteString(fmt.Sprintf("%s Fetching package details...", m.Spinner.View()))

	case ModalError:
		s.WriteString(errorStyle.Render("Failed to fetch/uninstall:\n"))
		s.WriteString(fmt.Sprintf("%v\n\nPress Esc to close.", m.Modal.Err))

	case ModalUninstalling:
		s.WriteString(fmt.Sprintf("%s Uninstalling %s...", m.Spinner.View(), m.Modal.Details.Name))

	case ModalReady:
		s.WriteString(titleStyle.Render(fmt.Sprintf("%s @ %s", m.Modal.Details.Name, m.Modal.Details.Version)) + "\n\n")

		if m.Modal.Details.Description != "" {
			s.WriteString(fmt.Sprintf("📝 %s\n\n", m.Modal.Details.Description))
		}
		if m.Modal.Details.Author != "" {
			s.WriteString(fmt.Sprintf("👤 Author: %s\n", m.Modal.Details.Author))
		}
		if m.Modal.Details.Homepage != "" {
			s.WriteString(fmt.Sprintf("🔗 URL: %s\n", m.Modal.Details.Homepage))
		}
		if m.Modal.Details.Size != "" {
			s.WriteString(fmt.Sprintf("📦 Size: %s\n", m.Modal.Details.Size))
		}
		if m.Modal.Details.Path != "" {
			s.WriteString(fmt.Sprintf("📂 Path: %s\n", m.Modal.Details.Path))
		}

		s.WriteString("\n" + strings.Repeat("─", 30) + "\n")
		s.WriteString(fmt.Sprintf("Press %s to Uninstall, or %s to close.", dangerStyle.Render("x"), selectedItemStyle.Render("Esc")))
	}

	modal := modalStyle.Render(s.String())

	// Simple centering logic over the background
	modalWidth := lipgloss.Width(modal)
	modalHeight := lipgloss.Height(modal)

	xOffset := (m.Width - modalWidth) / 2
	yOffset := (m.Height - modalHeight) / 2

	if xOffset < 0 {
		xOffset = 0
	}
	if yOffset < 0 {
		yOffset = 0
	}

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, modal)
}

func (m Model) viewLeftPane(width int) string {
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
		nameRender := style.Render(name)
		
		var versionRender string
		if lang.Version != "" {
			versionRender = versionStyle.Render(lang.Version)
		}

		if versionRender != "" {
			nameLen := lipgloss.Width(nameRender)
			verLen := lipgloss.Width(versionRender)
			
			padding := width - nameLen - verLen - 6
			if padding < 1 {
				padding = 1
			}
			padStr := strings.Repeat(".", padding)
			padRender := lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(padStr)
			
			s.WriteString(fmt.Sprintf("%s%s %s %s\n", cursor, nameRender, padRender, versionRender))
		} else {
			s.WriteString(fmt.Sprintf("%s%s\n", cursor, nameRender))
		}
	}
	return s.String()
}

func (m Model) viewMiddlePane(width int) string {
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
		nameRender := style.Render(name)
		
		var versionRender string
		if mgrData.Version != "" {
			versionRender = versionStyle.Render(mgrData.Version)
		}
		
		statusRender := ""
		if status != "" {
			statusRender = fmt.Sprintf(" [%s]", status)
		}
		
		if versionRender != "" {
			nameLen := lipgloss.Width(nameRender)
			verLen := lipgloss.Width(versionRender)
			statusLen := lipgloss.Width(statusRender)
			
			padding := width - nameLen - verLen - statusLen - 6
			if padding < 1 {
				padding = 1
			}
			padStr := strings.Repeat(".", padding)
			padRender := lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(padStr)
			
			s.WriteString(fmt.Sprintf("%s%s %s %s%s\n", cursor, nameRender, padRender, versionRender, statusRender))
		} else {
			s.WriteString(fmt.Sprintf("%s%s%s\n", cursor, nameRender, statusRender))
		}
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
