package tui

import "senk-tui/managers"

type DepsFetchedMsg struct {
	Index        int
	Version      string
	Dependencies []managers.Dependency
}

type FetchErrorMsg struct {
	Index int
	Err   error
}
