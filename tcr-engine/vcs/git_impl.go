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

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/go-git/go-billy/v5/helper/chroot"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/murex/tcr/tcr-engine/report"
	"path/filepath"
	"strings"
)

// GitImpl provides the implementation of the git interface
type GitImpl struct {
	baseDir                     string
	rootDir                     string
	remoteName                  string
	workingBranch               string
	workingBranchExistsOnRemote bool
	commitMessage               string
	pushEnabled                 bool
}

// New initializes the git implementation based on the provided directory from local clone
func New(dir string) (GitInterface, error) {
	var gitImpl = GitImpl{
		baseDir:       dir,
		remoteName:    DefaultRemoteName,
		commitMessage: DefaultCommitMessage,
		pushEnabled:   DefaultPushEnabled,
	}

	plainOpenOptions := git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: false,
	}
	repo, err := git.PlainOpenWithOptions(gitImpl.baseDir, &plainOpenOptions)
	if err != nil {
		return nil, err
	}
	r, _ := rootDir(repo)
	gitImpl.rootDir = filepath.Dir(r)

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	gitImpl.workingBranch = head.Name().Short()
	gitImpl.workingBranchExistsOnRemote, err = isBranchOnRemote(repo, gitImpl.workingBranch, gitImpl.remoteName)

	return &gitImpl, err
}

// isBranchOnRemote returns true is the provided branch exists on provided remote
func isBranchOnRemote(repo *git.Repository, branch, remote string) (bool, error) {
	remoteName := fmt.Sprintf("%v/%v", remote, branch)
	branches, err := remoteBranches(repo.Storer)
	if err != nil {
		return false, err
	}

	var found = false
	_ = branches.ForEach(func(branch *plumbing.Reference) error {
		found = found || strings.HasSuffix(branch.Name().Short(), remoteName)
		return nil
	})
	return found, nil
}

// remoteBranches returns the list of known remote branches
func remoteBranches(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()
	if err != nil {
		return nil, err
	}

	// We keep only remote branches, and ignore symbolic references
	return storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return ref.Name().IsRemote() && ref.Type() != plumbing.SymbolicReference
	}, refs), nil
}

// rootDir returns the local clone's root directory of provided repository
func rootDir(r *git.Repository) (string, error) {
	// Try to grab the repository Storer
	s, ok := r.Storer.(*filesystem.Storage)
	if !ok {
		return "", errors.New("repository storage is not filesystem.Storage")
	}

	// Try to get the underlying billy.Filesystem
	fs, ok := s.Filesystem().(*chroot.ChrootHelper)
	if !ok {
		return "", errors.New("filesystem is not chroot.ChrootHelper")
	}

	return fs.Root(), nil
}

// WorkingBranch returns the current git working branch
func (g *GitImpl) WorkingBranch() string {
	return g.workingBranch
}

// Commit restores to last commit.
// Current implementation uses a direct call to git
func (g *GitImpl) Commit() error {
	_ = traceGitCommand([]string{"commit", "--no-gpg-sign", "-am", g.commitMessage})
	// We ignore return code on purpose to prevent raising an error
	// when there is nothing to commit
	// TODO find a way to check beforehand if there is something to commit
	// ("git diff --exit-code --quiet HEAD" seems to do the trick)
	return nil
}

// Restore restores to last commit for the provided path.
// Current implementation uses a direct call to git
func (g *GitImpl) Restore(path string) error {
	report.PostWarning("Reverting ", path)
	return traceGitCommand([]string{"checkout", "HEAD", "--", path})
}

// Push runs a git push operation.
// Current implementation uses a direct call to git
func (g *GitImpl) Push() error {
	if !g.IsPushEnabled() {
		// There's nothing to do in this case
		return nil
	}

	report.PostInfo("Pushing changes to origin/", g.workingBranch)
	err := traceGitCommand([]string{"push", "--no-recurse-submodules", g.remoteName, g.workingBranch})
	if err == nil {
		g.workingBranchExistsOnRemote = true
	}
	return err
}

// Pull runs a git pull operation.
// Current implementation uses a direct call to git
func (g *GitImpl) Pull() error {
	if !g.workingBranchExistsOnRemote {
		report.PostInfo("Working locally on branch ", g.workingBranch)
		return nil
	}
	report.PostInfo("Pulling latest changes from ", g.remoteName, "/", g.workingBranch)
	return traceGitCommand([]string{"pull", "--no-recurse-submodules", g.remoteName, g.workingBranch})
}

// ListChanges returns the list of files modified since last commit
// Current implementation uses a direct call to git
func (g *GitImpl) ListChanges() (files []string, err error) {
	var gitOutput []byte
	gitOutput, err = runGitCommand([]string{"diff", "--name-only"})
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(gitOutput))
	for scanner.Scan() {
		files = append(files, filepath.Join(g.rootDir, scanner.Text()))
	}
	return
}

// EnablePush sets a flag allowing to turn on/off git push operations
func (g *GitImpl) EnablePush(flag bool) {
	if g.pushEnabled == flag {
		return
	}
	g.pushEnabled = flag
	autoPushStr := "off"
	if g.pushEnabled {
		autoPushStr = "on"
	}
	report.PostInfo(fmt.Sprintf("Git auto-push is turned %v", autoPushStr))
}

// IsPushEnabled indicates if git push operations are turned on
func (g *GitImpl) IsPushEnabled() bool {
	return g.pushEnabled
}

// traceGitCommand calls git command and reports its output
func traceGitCommand(params []string) error {
	output, err := runGitCommand(params)
	if len(output) > 0 {
		report.PostText(string(output))
	}
	return err
}

// runGitCommand calls git command in a separate process
func runGitCommand(params []string) (output []byte, err error) {
	return sh.Command("git", params).CombinedOutput()
}
