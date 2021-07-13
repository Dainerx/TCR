package tcr

import (
	"path/filepath"
	"strings"
)

type Language interface {
	name() string
	toolchain() string
	srcDirs() []string
	testDirs() []string
}

// ========================================================================

type JavaLanguage struct {
}

func (language JavaLanguage) name() string {
	return "java"
}

func (language JavaLanguage) toolchain() string {
	return "gradle"
}

func (language JavaLanguage) srcDirs() []string {
	return []string{
		filepath.Join("src", "main"),
	}
}

func (language JavaLanguage) testDirs() []string {
	return []string{
		filepath.Join("src", "test"),
	}
}

func (language JavaLanguage) matchesSrcFile(filename string) bool {
	extension := filepath.Ext(filename)
	switch strings.ToLower(extension) {
	case ".java":
		return true
	default:
		return false
	}
}

func (language JavaLanguage) matchesTestFile(filename string) bool {
	extension := filepath.Ext(filename)
	switch strings.ToLower(extension) {
	case ".java":
		return true
	default:
		return false
	}
}

// ========================================================================

type CppLanguage struct {
}

func (language CppLanguage) name() string {
	return "cpp"
}

func (language CppLanguage) toolchain() string {
	return "cmake"
}

func (language CppLanguage) srcDirs() []string {
	return []string{
		filepath.Join("src"),
		filepath.Join("include"),
	}
}

func (language CppLanguage) testDirs() []string {
	return []string{
		filepath.Join("test"),
	}
}

// ========================================================================

func dirsToWatch(baseDir string, lang Language) []string {
	dirList := append(lang.srcDirs(), lang.testDirs()...)
	for i := 0; i < len(dirList); i++ {
		dirList[i] = filepath.Join(baseDir, dirList[i])
	}
	//trace.Info(dirList)
	return dirList
}
