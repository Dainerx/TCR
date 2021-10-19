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

package cli

import (
	"github.com/murex/tcr-engine/engine"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_confirm_with_default_answer_to_yes(t *testing.T) {
	assertConfirmBehaviour(t, []byte{enterKey}, true, true)
}

func Test_confirm_with_default_answer_to_no(t *testing.T) {
	assertConfirmBehaviour(t, []byte{enterKey}, false, false)
}

func Test_confirm_with_a_yes_answer(t *testing.T) {
	assertConfirmBehaviour(t, []byte{'y'}, false, true)
	assertConfirmBehaviour(t, []byte{'Y'}, false, true)
}

func Test_confirm_with_a_no_answer(t *testing.T) {
	assertConfirmBehaviour(t, []byte{'n'}, true, false)
	assertConfirmBehaviour(t, []byte{'N'}, true, false)
}

func Test_confirm_question_with_default_answer_to_no(t *testing.T) {
	assert.Equal(t, "[y/N]", yesOrNoAdvice(false))
}

func Test_confirm_question_with_default_answer_to_yes(t *testing.T) {
	assert.Equal(t, "[Y/n]", yesOrNoAdvice(true))
}

func assertConfirmBehaviour(t *testing.T, input []byte, defaultValue bool, expected bool) {
	stdin := os.Stdin
	stdout := os.Stdout
	// Restore stdin and stdout right after the test.
	defer func() { os.Stdin = stdin; os.Stdout = stdout }()
	// We mock stdin so that we can simulate a key press
	os.Stdin = mockStdin(t, input)
	// Displayed info on stdout is useless for the test
	os.Stdout = os.NewFile(0, os.DevNull)

	term := New(engine.Params{})
	assert.Equal(t, expected, term.Confirm("", defaultValue))
}

func mockStdin(t *testing.T, input []byte) *os.File {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}
	_ = w.Close()
	return r
}
