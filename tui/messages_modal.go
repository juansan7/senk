package tui

import "senk-tui/managers"

type DetailsFetchedMsg struct {
	Details managers.PackageDetails
}

type DetailsErrorMsg struct {
	Err error
}

type UninstallCompleteMsg struct{}

type UninstallErrorMsg struct {
	Err error
}
