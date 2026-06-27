package tui

import "senk-tui/managers"

type ModalState int

const (
	ModalClosed ModalState = iota
	ModalLoading
	ModalReady
	ModalUninstalling
	ModalError
)

type ModalData struct {
	State   ModalState
	Details managers.PackageDetails
	Err     error
	// To know what to refresh after uninstall
	LangIndex int
	MgrIndex  int
}
