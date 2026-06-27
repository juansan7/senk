package tui

import (
	"senk-tui/managers"

	"github.com/charmbracelet/bubbles/spinner"
)

type ManagerState int

const (
	StateLoading ManagerState = iota
	StateDone
	StateError
)

type ManagerData struct {
	Manager      managers.PackageManager
	State        ManagerState
	Dependencies []managers.Dependency
	Err          error
}

type FocusState int

const (
	FocusLeft FocusState = iota
	FocusRight
)

type Model struct {
	Managers   []*ManagerData
	LeftIndex  int
	RightIndex int
	Focus      FocusState
	Width      int
	Height     int
	Spinner    spinner.Model
}

func NewModel(availableManagers []managers.PackageManager) Model {
	var mData []*ManagerData
	for _, m := range availableManagers {
		mData = append(mData, &ManagerData{
			Manager: m,
			State:   StateLoading,
		})
	}

	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		Managers: mData,
		Focus:    FocusLeft,
		Spinner:  s,
	}
}
