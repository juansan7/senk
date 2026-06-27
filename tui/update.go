package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"senk-tui/managers"
)

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.Spinner.Tick)

	for lIdx, lang := range m.Languages {
		cmds = append(cmds, fetchLangVersionCmd(lIdx, lang.Language))
		for mIdx, mgr := range lang.Managers {
			cmds = append(cmds, fetchMgrCmd(lIdx, mIdx, mgr.Manager))
		}
	}

	return tea.Batch(cmds...)
}

func fetchLangVersionCmd(lIdx int, lang managers.Language) tea.Cmd {
	return func() tea.Msg {
		version, _ := lang.LanguageVersion()
		return LangVersionFetchedMsg{LangIndex: lIdx, Version: version}
	}
}

func fetchMgrCmd(lIdx int, mIdx int, manager managers.PackageManager) tea.Cmd {
	return func() tea.Msg {
		version, _ := manager.ManagerVersion()
		deps, err := manager.Fetch()
		if err != nil {
			return FetchErrorMsg{LangIndex: lIdx, MgrIndex: mIdx, Err: err}
		}
		return DepsFetchedMsg{LangIndex: lIdx, MgrIndex: mIdx, Version: version, Dependencies: deps}
	}
}

func fetchDetailsCmd(manager managers.PackageManager, dep managers.Dependency) tea.Cmd {
	return func() tea.Msg {
		details, err := manager.FetchDetails(dep)
		if err != nil {
			return DetailsErrorMsg{Err: err}
		}
		return DetailsFetchedMsg{Details: details}
	}
}

func uninstallCmd(manager managers.PackageManager, dep managers.Dependency) tea.Cmd {
	return func() tea.Msg {
		err := manager.Uninstall(dep)
		if err != nil {
			return UninstallErrorMsg{Err: err}
		}
		return UninstallCompleteMsg{}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		// Modal Interception
		if m.Modal.State != ModalClosed {
			switch msg.String() {
			case "esc", "q":
				if m.Modal.State != ModalUninstalling {
					m.Modal.State = ModalClosed
				}
			case "x", "backspace":
				if m.Modal.State == ModalReady {
					m.Modal.State = ModalUninstalling
					activeManager := m.Languages[m.Modal.LangIndex].Managers[m.Modal.MgrIndex].Manager
					dep := managers.Dependency{Name: m.Modal.Details.Name, Version: m.Modal.Details.Version}
					return m, uninstallCmd(activeManager, dep)
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab", "l", "right", "enter":
			if m.Focus == FocusLeft && len(m.Languages) > 0 {
				m.Focus = FocusMiddle
				m.MgrIndex = 0
				m.DepIndex = 0
			} else if m.Focus == FocusMiddle && len(m.Languages) > 0 && len(m.Languages[m.LangIndex].Managers) > 0 {
				m.Focus = FocusRight
				m.DepIndex = 0
			} else if m.Focus == FocusRight && msg.String() == "enter" {
				// Open Modal
				activeManager := m.Languages[m.LangIndex].Managers[m.MgrIndex]
				if len(activeManager.Dependencies) > 0 {
					dep := activeManager.Dependencies[m.DepIndex]
					m.Modal.State = ModalLoading
					m.Modal.LangIndex = m.LangIndex
					m.Modal.MgrIndex = m.MgrIndex
					return m, fetchDetailsCmd(activeManager.Manager, dep)
				}
			}

		case "esc", "left", "h":
			if m.Focus == FocusRight {
				m.Focus = FocusMiddle
			} else if m.Focus == FocusMiddle {
				m.Focus = FocusLeft
			}

		case "up", "k":
			if m.Focus == FocusLeft && m.LangIndex > 0 {
				m.LangIndex--
				m.MgrIndex = 0
				m.DepIndex = 0
			} else if m.Focus == FocusMiddle && m.MgrIndex > 0 {
				m.MgrIndex--
				m.DepIndex = 0
			} else if m.Focus == FocusRight && m.DepIndex > 0 {
				m.DepIndex--
			}

		case "down", "j":
			if m.Focus == FocusLeft && m.LangIndex < len(m.Languages)-1 {
				m.LangIndex++
				m.MgrIndex = 0
				m.DepIndex = 0
			} else if m.Focus == FocusMiddle && m.MgrIndex < len(m.Languages[m.LangIndex].Managers)-1 {
				m.MgrIndex++
				m.DepIndex = 0
			} else if m.Focus == FocusRight {
				if len(m.Languages) > 0 && len(m.Languages[m.LangIndex].Managers) > 0 {
					depsCount := len(m.Languages[m.LangIndex].Managers[m.MgrIndex].Dependencies)
					if m.DepIndex < depsCount-1 {
						m.DepIndex++
					}
				}
			}
		}

	case DetailsFetchedMsg:
		m.Modal.State = ModalReady
		m.Modal.Details = msg.Details

	case DetailsErrorMsg:
		m.Modal.State = ModalError
		m.Modal.Err = msg.Err

	case UninstallCompleteMsg:
		m.Modal.State = ModalClosed
		// Re-fetch the manager to refresh the list
		activeManager := m.Languages[m.Modal.LangIndex].Managers[m.Modal.MgrIndex]
		activeManager.State = StateLoading
		cmds = append(cmds, fetchMgrCmd(m.Modal.LangIndex, m.Modal.MgrIndex, activeManager.Manager))

	case UninstallErrorMsg:
		m.Modal.State = ModalError
		m.Modal.Err = msg.Err

	case LangVersionFetchedMsg:
		if msg.LangIndex >= 0 && msg.LangIndex < len(m.Languages) {
			m.Languages[msg.LangIndex].Version = msg.Version
		}

	case DepsFetchedMsg:
		if msg.LangIndex >= 0 && msg.LangIndex < len(m.Languages) {
			if msg.MgrIndex >= 0 && msg.MgrIndex < len(m.Languages[msg.LangIndex].Managers) {
				m.Languages[msg.LangIndex].Managers[msg.MgrIndex].Version = msg.Version
				m.Languages[msg.LangIndex].Managers[msg.MgrIndex].Dependencies = msg.Dependencies
				m.Languages[msg.LangIndex].Managers[msg.MgrIndex].State = StateDone
			}
		}

	case FetchErrorMsg:
		if msg.LangIndex >= 0 && msg.LangIndex < len(m.Languages) {
			if msg.MgrIndex >= 0 && msg.MgrIndex < len(m.Languages[msg.LangIndex].Managers) {
				m.Languages[msg.LangIndex].Managers[msg.MgrIndex].Err = msg.Err
				m.Languages[msg.LangIndex].Managers[msg.MgrIndex].State = StateError
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}
	return m, nil
}
