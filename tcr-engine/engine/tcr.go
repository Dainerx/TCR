/*
Copyright (c) 2021 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package engine

import (
	"github.com/murex/tcr-engine/filesystem"
	"github.com/murex/tcr-engine/language"
	"github.com/murex/tcr-engine/report"
	"github.com/murex/tcr-engine/role"
	"github.com/murex/tcr-engine/runmode"
	"github.com/murex/tcr-engine/settings"
	"github.com/murex/tcr-engine/timer"
	"github.com/murex/tcr-engine/toolchain"
	"github.com/murex/tcr-engine/ui"
	"github.com/murex/tcr-engine/vcs"
	"gopkg.in/tomb.v2"
	"os"
	"path/filepath"
	"time"
)

var (
	mode            runmode.RunMode
	uitf            ui.UserInterface
	git             vcs.GitInterface
	lang            language.Language
	tchn            toolchain.Toolchain
	sourceTree      filesystem.SourceTree
	pollingPeriod   time.Duration
	mobTurnDuration time.Duration
)

// Init initializes the TCR engine with the provided parameters, and wires it to the user interface.
// This function should be called only once during the lifespan of the application
func Init(u ui.UserInterface, params Params) {
	var err error

	uitf = u

	report.PostInfo("Starting TCR version ", Version, "...")

	SetRunMode(params.Mode)
	pollingPeriod = params.PollingPeriod

	sourceTree, err = filesystem.New(params.BaseDir)
	handleError(err)
	report.PostInfo("Working directory is ", sourceTree.GetBaseDir())
	lang, err = language.DetectLanguage(sourceTree.GetBaseDir())
	handleError(err)
	tchn, err = toolchain.New(params.Toolchain, lang)
	handleError(err)
	git, err = vcs.New(sourceTree.GetBaseDir())
	handleError(err)
	git.EnablePush(params.AutoPush)

	if settings.EnableMobTimer {
		mobTurnDuration = params.MobTurnDuration
		report.PostInfo("Mob timer duration is ", mobTurnDuration)
	}

	uitf.ShowRunningMode(mode)
	uitf.ShowSessionInfo()
	warnIfOnRootBranch(git.WorkingBranch())
}

func warnIfOnRootBranch(branch string) {
	for _, b := range []string{"main", "master"} {
		if b == branch {
			if !uitf.Confirm("Running TCR on branch \""+branch+"\" is not recommended", false) {
				Quit()
			}
			break
		}
	}
}

// ToggleAutoPush toggles git auto-push state
func ToggleAutoPush() {
	git.EnablePush(!git.IsPushEnabled())
}

// SetAutoPush sets git auto-push to the provided value
func SetAutoPush(ap bool) {
	git.EnablePush(ap)
}

// RunAsDriver tells TCR engine to start running with driver role
func RunAsDriver() {
	var mobTimer *timer.PeriodicReminder
	if settings.EnableMobTimer {
		mobTimer = timer.NewMobTurnCountdown(mode, mobTurnDuration)
	}

	go fromBirthTillDeath(
		func() {
			uitf.NotifyRoleStarting(role.Driver{})
			_ = git.Pull()
			if settings.EnableMobTimer {
				mobTimer.Start()
			}
		},
		func(interrupt <-chan bool) bool {
			inactivityTeaser := timer.GetInactivityTeaserInstance()
			inactivityTeaser.Start()
			if waitForChange(interrupt) {
				// Some file changes were detected
				inactivityTeaser.Reset()
				runTCR()
				inactivityTeaser.Start()
				return true
			}
			// If we arrive here this means that the end of waitForChange
			// was triggered by the user
			inactivityTeaser.Reset()
			return false
		},
		func() {
			if settings.EnableMobTimer {
				mobTimer.Stop()
			}
			uitf.NotifyRoleEnding(role.Driver{})
		},
	)
}

// RunAsNavigator tells TCR engine to start running with navigator role
func RunAsNavigator() {
	go fromBirthTillDeath(
		func() {
			uitf.NotifyRoleStarting(role.Navigator{})
		},
		func(interrupt <-chan bool) bool {
			select {
			case <-interrupt:
				return false
			default:
				_ = git.Pull()
				time.Sleep(pollingPeriod)
				return true
			}
		},
		func() {
			uitf.NotifyRoleEnding(role.Navigator{})
		},
	)
}

// shoot channel is used to handle interruptions coming from the UI
var shoot chan bool

// Stop is the entry point for telling TCR engine to stop its current operations
func Stop() {
	shoot <- true
}

func fromBirthTillDeath(
	birth func(),
	dailyLife func(interrupt <-chan bool) bool,
	death func(),
) {
	var tmb tomb.Tomb
	shoot = make(chan bool)

	// The goroutine doing the work
	tmb.Go(func() error {
		birth()
		for oneMoreDay := true; oneMoreDay; {
			oneMoreDay = dailyLife(shoot)
		}
		death()
		return nil
	})

	err := tmb.Wait()
	if err != nil {
		report.PostError("tmb.Wait(): ", err)
	}
}

func waitForChange(interrupt <-chan bool) bool {
	report.PostInfo("Going to sleep until something interesting happens")
	return sourceTree.Watch(
		language.DirsToWatch(sourceTree.GetBaseDir(), lang),
		lang.IsSrcFile,
		interrupt)
}

func runTCR() {
	if build() != nil {
		return
	}
	if test() == nil {
		commit()
	} else {
		revert()
	}
}

func build() error {
	report.PostInfo("Launching Build")
	err := tchn.RunBuild()
	if err != nil {
		report.PostWarning("There are build errors! I can't go any further")
	}
	return err
}

func test() error {
	report.PostInfo("Running Tests")
	err := tchn.RunTests()
	if err != nil {
		report.PostWarning("Some tests are failing! That's unfortunate")
	}
	return err
}

func commit() {
	report.PostInfo("Committing changes on branch ", git.WorkingBranch())
	_ = git.Commit()
	_ = git.Push()
}

func revert() {
	report.PostWarning("Reverting changes")
	for _, dir := range lang.SrcDirs() {
		_ = git.Restore(filepath.Join(sourceTree.GetBaseDir(), dir))
	}
}

// GetSessionInfo provides the information (as strings) related to the current TCR session.
// Used mainly by the user interface packages to retrieve and display this information
func GetSessionInfo() (d string, l string, t string, ap bool, b string) {
	d = sourceTree.GetBaseDir()
	l = lang.Name()
	t = tchn.Name()
	ap = git.IsPushEnabled()
	b = git.WorkingBranch()

	return d, l, t, ap, b
}

// SetRunMode sets the run mode for TCR engine
func SetRunMode(m runmode.RunMode) {
	mode = m
}

// Quit is the exit point for TCR application
func Quit() {
	report.PostInfo("That's All Folks!")
	time.Sleep(1 * time.Millisecond)
	os.Exit(0)
}

func handleError(err error) {
	if err != nil {
		report.PostError(err)
		time.Sleep(1 * time.Millisecond)
		os.Exit(1)
	}
}
