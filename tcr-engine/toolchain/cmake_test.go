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

package toolchain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_cmake_is_a_built_in_toolchain(t *testing.T) {
	assert.True(t, isBuiltIn("cmake"))
}

func Test_cmake_toolchain_is_supported(t *testing.T) {
	assert.True(t, isSupported("cmake"))
}

func Test_cmake_toolchain_name_is_case_insensitive(t *testing.T) {
	assert.True(t, isSupported("cmake"))
	assert.True(t, isSupported("Cmake"))
	assert.True(t, isSupported("CMAKE"))
}

func Test_cmake_toolchain_initialization(t *testing.T) {
	toolchain, err := Get("cmake")
	assert.Equal(t, "cmake", toolchain.GetName())
	assert.Zero(t, err)
}

func Test_cmake_toolchain_name(t *testing.T) {
	toolchain, _ := Get("cmake")
	assert.Equal(t, "cmake", toolchain.GetName())
}

func Test_cmake_toolchain_build_command_args(t *testing.T) {
	toolchain, _ := Get("cmake")
	assert.Equal(t, []string{
		"--build", "build",
		"--config", "Debug",
	}, toolchain.BuildCommandArgs())
}

func Test_cmake_toolchain_returns_error_when_build_fails(t *testing.T) {
	// Note: this passes not due to cmake return value, but due to absence of cmake command
	toolchain, _ := Get("cmake")
	runFromDir(t, testDataRootDir,
		func(t *testing.T) {
			assert.NotZero(t, toolchain.RunBuild())
		})
}

// TODO Figure out a way to provide a cmake wrapper
//func test_cmake_toolchain_returns_ok_when_build_passes(t *testing.T) {
//  toolchain, _ := Get("cmake")
//	runFromDir(t, testDataDirCpp,
//		func(t *testing.T) {
//			assert.Zero(t, toolchain.RunBuild())
//		})
//}

func Test_cmake_toolchain_test_command_args(t *testing.T) {
	toolchain, _ := Get("cmake")
	assert.Equal(t, []string{
		"--output-on-failure",
		"--test-dir", "build",
		"--build-config", "Debug",
	}, toolchain.TestCommandArgs())
}

func Test_cmake_toolchain_returns_error_when_tests_fail(t *testing.T) {
	// Note: this passes not due to ctest return value, but due to absence of ctest command
	toolchain, _ := Get("cmake")
	runFromDir(t, testDataRootDir,
		func(t *testing.T) {
			assert.NotZero(t, toolchain.RunTests())
		})
}

// TODO Figure out a way to provide a cmake wrapper
//func Test_cmake_toolchain_returns_ok_when_tests_pass(t *testing.T) {
//  toolchain, _ := Get("cmake")
//	runFromDir(t, testDataDirCpp,
//		func(t *testing.T) {
//			assert.Zero(t, toolchain.RunTests())
//		})
//}

func Test_cmake_toolchain_supported_platforms(t *testing.T) {
	// Cf. https://cmake.org/download/ for list of cmake supported platforms
	toolchain, _ := Get("cmake")

	// Windows platforms
	assert.True(t, toolchain.supportsPlatform(OsWindows, Arch386))
	assert.True(t, toolchain.supportsPlatform(OsWindows, ArchAmd64))
	assert.False(t, toolchain.supportsPlatform(OsWindows, ArchArm64))

	// Darwin platforms
	assert.False(t, toolchain.supportsPlatform(OsDarwin, Arch386))
	assert.True(t, toolchain.supportsPlatform(OsDarwin, ArchAmd64))
	assert.True(t, toolchain.supportsPlatform(OsDarwin, ArchArm64))

	// Linux platforms
	assert.False(t, toolchain.supportsPlatform(OsLinux, Arch386))
	assert.True(t, toolchain.supportsPlatform(OsLinux, ArchAmd64))
	assert.True(t, toolchain.supportsPlatform(OsLinux, ArchArm64))
}
