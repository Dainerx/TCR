package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/mengdaming/tcr/tcr/engine"
	"github.com/mengdaming/tcr/tcr/report"
	"github.com/mengdaming/tcr/tcr/role"
	"github.com/mengdaming/tcr/tcr/runmode"
	"github.com/mengdaming/tcr/tcr/ui"
	"github.com/mengdaming/tcr/tcr/ui/cli"
	"image/color"
	"strings"
)

type GUI struct {
	app                  fyne.App
	win                  fyne.Window
	directoryLabel       *widget.Label
	languageLabel        *widget.Label
	toolchainLabel       *widget.Label
	branchLabel          *widget.Label
	autoPushToggle       *widget.Check
	startNavigatorButton *widget.Button
	startDriverButton    *widget.Button
	stopButton           *widget.Button
	traceVBox            *fyne.Container
	traceArea            *container.Scroll
}

var (
	redColor    = color.RGBA{R: 255, G: 0, B: 0}
	cyanColor   = color.RGBA{R: 0, G: 139, B: 139}
	yellowColor = color.RGBA{R: 255, G: 255, B: 0}
	whiteColor  = color.RGBA{R: 255, G: 255, B: 255}
)

// TODO Remove once all GUI implementations are available
var term ui.UserInterface

func New() ui.UserInterface {
	term = cli.New()
	var gui = GUI{}
	gui.initApp()
	gui.startReporting()

	return &gui
}

func (gui *GUI) startReporting() chan bool {
	// TODO Unsubscribe on quit
	return report.Subscribe(func(msg report.Message) {
		switch msg.Type {
		case report.Normal:
			gui.trace(msg.Text)
		case report.Title:
			gui.title(msg.Text)
		case report.Info:
			gui.info(msg.Text)
		case report.Warning:
			gui.warning(msg.Text)
		case report.Error:
			gui.error(msg.Text)
		}
	})
}

func (gui *GUI) RunInMode(mode runmode.RunMode) {
	// TODO setup according to mode value
	gui.win.ShowAndRun()
}

func (gui *GUI) ShowRunningMode(mode runmode.RunMode) {
	// TODO Replace with GUI-specific implementation
	term.ShowRunningMode(mode)
}

func (gui *GUI) NotifyRoleStarting(r role.Role) {
	report.PostTitle("Starting as a ", strings.Title(r.Name()))
}

func (gui *GUI) NotifyRoleEnding(r role.Role) {
	report.PostInfo("Leaving ", strings.Title(r.Name()), " role")
}

func (gui *GUI) ShowSessionInfo() {
	d, l, t, ap, b := engine.GetSessionInfo()

	gui.directoryLabel.SetText(fmt.Sprintf("Directory: %v", d))
	gui.languageLabel.SetText(fmt.Sprintf("Language: %v", l))
	gui.toolchainLabel.SetText(fmt.Sprintf("Toolchain: %v", t))
	gui.branchLabel.SetText(fmt.Sprintf("Branch: %v", b))
	gui.autoPushToggle.SetChecked(ap)
}

func (gui *GUI) info(a ...interface{}) {
	gui.traceText(cyanColor, a...)
}

func (gui *GUI) title(a ...interface{}) {
	gui.traceLine()
	gui.traceVBox.Add(widget.NewLabelWithStyle(
		fmt.Sprint(a...),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	))
	gui.traceLine()
	gui.scrollTraceToBottom()
}

func (gui *GUI) warning(a ...interface{}) {
	gui.traceText(yellowColor, a...)
}

func (gui *GUI) error(a ...interface{}) {
	gui.traceText(redColor, a...)
}

func (gui *GUI) trace(a ...interface{}) {
	gui.traceText(whiteColor, a...)
}

func (gui *GUI) traceText(col color.Color, a ...interface{}) {
	for _, s := range strings.Split(fmt.Sprint(a...), "\n") {
		gui.traceVBox.Add(canvas.NewText(s, col))
	}
	gui.scrollTraceToBottom()
}

func (gui *GUI) traceLine() {
	gui.traceVBox.Add(widget.NewSeparator())
	gui.scrollTraceToBottom()
}

func (gui *GUI) scrollTraceToBottom() {
	// The ScrollToTop() call below is some kind of workaround to ensure
	// that the UI indeed refreshes and scrolls to bottom when ScrollToBottom() is called
	gui.traceArea.ScrollToTop()
	gui.traceArea.ScrollToBottom()
}

func (gui *GUI) Confirm(message string, def bool) bool {
	// TODO Replace with GUI-specific implementation
	return term.Confirm(message, def)
}

func (gui *GUI) initApp() {
	gui.app = app.New()
	// TODO Add a TCR-Specific icon
	gui.app.SetIcon(theme.FyneLogo())
	gui.win = gui.app.NewWindow("TCR")
	gui.win.Resize(fyne.NewSize(400, 800))

	// TODO Refactor into smaller functions, one for each main UI area

	// Action Buttons container

	gui.startDriverButton = widget.NewButtonWithIcon("Start as Driver", theme.MediaPlayIcon(), func() {
		// TODO Remove once everything works as expected
		//report.PostWarning("Start as Driver Pushed")
		gui.startDriverButton.Disable()
		gui.startNavigatorButton.Disable()
		gui.stopButton.Enable()
		engine.RunAsDriver()
	})
	gui.startNavigatorButton = widget.NewButtonWithIcon("Start as Navigator", theme.MediaPlayIcon(), func() {
		// TODO Remove once everything works as expected
		//report.PostWarning("Start as Navigator Pushed")
		gui.startDriverButton.Disable()
		gui.startNavigatorButton.Disable()
		gui.stopButton.Enable()

		engine.RunAsNavigator()
	})
	gui.stopButton = widget.NewButtonWithIcon("Stop", theme.MediaStopIcon(), func() {
		// TODO Remove once everything works as expected
		//report.PostWarning("Stop Pushed")
		gui.stopButton.Disable()
		gui.startDriverButton.Enable()
		gui.startNavigatorButton.Enable()
		engine.Stop()
	})
	actionBar := container.NewHBox(
		layout.NewSpacer(),
		gui.startDriverButton,
		gui.startNavigatorButton,
		layout.NewSpacer(),
		gui.stopButton,
		layout.NewSpacer(),
	)

	// Initial state
	// TODO encapsulate in a single function to enforce consistency
	gui.startDriverButton.Enable()
	gui.startNavigatorButton.Enable()
	gui.stopButton.Disable()

	// Trace container

	gui.traceVBox = container.NewVBox(
		widget.NewLabelWithStyle("Welcome to TCR!",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		))
	gui.traceArea = container.NewVScroll(
		gui.traceVBox,
	)

	// Session Information container

	gui.directoryLabel = widget.NewLabel("Directory")
	gui.languageLabel = widget.NewLabel("Language")
	gui.toolchainLabel = widget.NewLabel("Toolchain")
	gui.branchLabel = widget.NewLabel("Branch")
	gui.autoPushToggle = widget.NewCheck("Auto-Push", func(checked bool) {
		engine.ToggleAutoPush()
		autoPushStr := "off"
		if checked {
			autoPushStr = "on"
		}
		report.PostInfo(fmt.Sprintf("Git auto-push is turned %v", autoPushStr))
	})
	sessionInfo := container.NewVBox(
		container.NewHBox(
			gui.directoryLabel,
		),
		widget.NewSeparator(),
		container.NewHBox(
			gui.languageLabel,
			widget.NewSeparator(),
			gui.toolchainLabel,
			widget.NewSeparator(),
			gui.branchLabel,
			widget.NewSeparator(),
			gui.autoPushToggle,
		),
		widget.NewSeparator(),
	)

	// Top level container

	topLevel := container.New(layout.NewBorderLayout(
		sessionInfo, actionBar, nil, nil),
		sessionInfo, actionBar, gui.traceArea)

	gui.win.SetContent(topLevel)
}
