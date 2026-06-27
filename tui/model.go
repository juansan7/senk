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
	Version      string
	State        ManagerState
	Dependencies []managers.Dependency
	Err          error
}

type LanguageData struct {
	Language managers.Language
	Version  string
	Managers []*ManagerData
}

type FocusState int

const (
	FocusLeft FocusState = iota
	FocusMiddle
	FocusRight
)

type Model struct {
	Languages []*LanguageData
	LangIndex int
	MgrIndex  int
	DepIndex  int
	Focus     FocusState
	Width     int
	Height    int
	Spinner   spinner.Model
}

func NewModel(availableLanguages []managers.Language) Model {
	var lData []*LanguageData

	for _, l := range availableLanguages {
		var mData []*ManagerData
		for _, m := range l.Managers() {
			if m.IsInstalled() {
				mData = append(mData, &ManagerData{
					Manager: m,
					State:   StateLoading,
				})
			}
		}

		if len(mData) > 0 {
			lData = append(lData, &LanguageData{
				Language: l,
				Managers: mData,
			})
		}
	}

	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		Languages: lData,
		Focus:     FocusLeft,
		Spinner:   s,
	}
}
