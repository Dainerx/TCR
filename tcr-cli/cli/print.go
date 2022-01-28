/*
Copyright (c) 2022 Murex

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

package cli

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/logrusorgru/aurora"
	"strconv"
	"strings"
)

const (
	horizontalLineCharacter = "-" // Character used for printing horizontal lines
	defaultTerminalWidth    = 80  // Default terminal width if current terminal is not recognized
)

var (
	colorizer       = aurora.NewAurora(true)
	linePrefix      = ""
	tputCmdDisabled = false
)

func setLinePrefix(value string) {
	linePrefix = value
}

func printPrefixedAndColored(fgColor aurora.Color, message string) {
	setupConsole()
	fmt.Println(
		colorizer.Colorize(linePrefix, fgColor),
		colorizer.Colorize(message, fgColor))
}

func printInCyan(a ...interface{}) {
	printPrefixedAndColored(aurora.CyanFg, fmt.Sprint(a...))
}

func printInYellow(a ...interface{}) {
	printPrefixedAndColored(aurora.YellowFg, fmt.Sprint(a...))
}

func printInGreen(a ...interface{}) {
	printPrefixedAndColored(aurora.GreenFg, fmt.Sprint(a...))
}

func printInRed(a ...interface{}) {
	printPrefixedAndColored(aurora.RedFg, fmt.Sprint(a...))
}

func printUntouched(a ...interface{}) {
	fmt.Println(a...)
}

func printHorizontalLine() {
	lineWidth := getTerminalColumns() - len(linePrefix) - 2
	if lineWidth < 0 {
		lineWidth = 0
	}
	printInCyan(strings.Repeat(horizontalLineCharacter, lineWidth))
}

// getTerminalColumns returns the terminal's current number of column. If anything goes wrong (for
// example when running from Windows PowerShell), we fallback on a fixed number of columns
func getTerminalColumns() int {
	if tputCmdDisabled {
		return defaultTerminalWidth
	}
	output, err := sh.Command("tput", "cols").Output()
	if err != nil {
		tputCmdDisabled = true
		return defaultTerminalWidth
	}
	columns, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		tputCmdDisabled = true
		return defaultTerminalWidth
	}
	return columns
}
