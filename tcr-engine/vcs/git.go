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

package vcs

const (
	// DefaultRemoteName is the alias used by default for the git remote repository
	DefaultRemoteName = "origin"
	// DefaultCommitMessage is the message used by default by TCR every time it does a git commit
	DefaultCommitMessage = "TCR"
	// DefaultPushEnabled indicates the default state for git auto-push option
	DefaultPushEnabled = false
)

// GitInterface provides the interface that a git implementation must satisfy for TCR engine to be
// able to interact with git
type GitInterface interface {
	WorkingBranch() string
	Commit() error
	Restore(dir string) error
	Push() error
	Pull() error
	EnablePush(flag bool)
	IsPushEnabled() bool
}
