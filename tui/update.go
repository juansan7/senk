package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"senk-tui/managers"
)

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.Spinner.Tick)

	for i := range m.Managers {
		cmds = append(cmds, fetchCmd(i, m.Managers[i].Manager))
	}

	return tea.Batch(cmds...)
}

func fetchCmd(index int, manager managers.PackageManager) tea.Cmd {
	return func() tea.Msg {
		version, _ := manager.ManagerVersion()
		deps, err := manager.Fetch()
		if err != nil {
			return FetchErrorMsg{Index: index, Err: err}
		}
		return DepsFetchedMsg{Index: index, Version: version, Dependencies: deps}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab", "right", "l", "enter":
			if m.Focus == FocusLeft && len(m.Managers) > 0 {
				m.Focus = FocusRight
				m.RightIndex = 0 // Reset right pane scroll when entering
			}
		case "esc", "left", "h":
			if m.Focus == FocusRight {
				m.Focus = FocusLeft
			}
		case "up", "k":
			if m.Focus == FocusLeft {
				if m.LeftIndex > 0 {
					m.LeftIndex--
				}
			} else if m.Focus == FocusRight {
				if m.RightIndex > 0 {
					m.RightIndex--
				}
			}
		case "down", "j":
			if m.Focus == FocusLeft {
				if m.LeftIndex < len(m.Managers)-1 {
					m.LeftIndex++
				}
			} else if m.Focus == FocusRight {
				if len(m.Managers) > 0 {
					depsCount := len(m.Managers[m.LeftIndex].Dependencies)
					if m.RightIndex < depsCount-1 {
						m.RightIndex++
					}
				}
			}
		}

	case DepsFetchedMsg:
		if msg.Index >= 0 && msg.Index < len(m.Managers) {
			m.Managers[msg.Index].Version = msg.Version
			m.Managers[msg.Index].Dependencies = msg.Dependencies
			m.Managers[msg.Index].State = StateDone
		}
	case FetchErrorMsg:
		if msg.Index >= 0 && msg.Index < len(m.Managers) {
			m.Managers[msg.Index].Err = msg.Err
			m.Managers[msg.Index].State = StateError
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}
