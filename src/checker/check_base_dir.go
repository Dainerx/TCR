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

package checker

import (
	"github.com/murex/tcr/params"
)

func checkBaseDirectory(p params.Params) (cr *CheckResults) {
	cr = NewCheckResults("base directory")

	if p.BaseDir == "" {
		cr.ok("base directory parameter is not set explicitly")
	} else {
		cr.ok("base directory parameter is ", p.BaseDir)
	}

	if checkEnv.sourceTreeErr != nil {
		cr.add(checkpointsForDirAccessError(p.BaseDir, checkEnv.sourceTreeErr)...)
	} else {
		cr.ok("base directory absolute path is ", checkEnv.sourceTree.GetBaseDir())
	}
	return cr
}
