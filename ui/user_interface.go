package ui

import (
	"github.com/mengdaming/tcr/role"
	"github.com/mengdaming/tcr/runmode"
)

type UserInterface interface {
	Start(mode runmode.RunMode)
	ShowRunningMode(mode runmode.RunMode)
	NotifyRoleStarting(r role.Role)
	NotifyRoleEnding(r role.Role)
	ShowSessionInfo()
	Confirm(message string, def bool) bool
	StartReporting()
	StopReporting()
}