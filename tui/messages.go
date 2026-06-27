package tui

import "senk-tui/managers"

type DepsFetchedMsg struct {
	Index        int
	Dependencies []managers.Dependency
}

type FetchErrorMsg struct {
	Index int
	Err   error
}
