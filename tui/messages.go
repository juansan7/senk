package tui

import "senk-tui/managers"

type LangVersionFetchedMsg struct {
	LangIndex int
	Version   string
}

type DepsFetchedMsg struct {
	LangIndex    int
	MgrIndex     int
	Version      string
	Dependencies []managers.Dependency
}

type FetchErrorMsg struct {
	LangIndex int
	MgrIndex  int
	Err       error
}
