package tcr

import (
	"github.com/codeskyblue/go-sh"
	"github.com/mengdaming/tcr/trace"
	"os"
	"path/filepath"
)

type Toolchain interface {
	name() string
	runBuild() error
	runTests() error
	buildCommandName() string
	buildCommandArgs() []string
	testCommandName() string
	testCommandArgs() []string
	supports(language Language) bool
}

func NewToolchain(name string, language Language) Toolchain {
	var toolchain Toolchain = nil
	switch name {
	case GradleToolchain{}.name():
		toolchain = GradleToolchain{}
	case MavenToolchain{}.name():
		toolchain = MavenToolchain{}
	case CmakeToolchain{}.name():
		toolchain = CmakeToolchain{}
	case "":
		toolchain = defaultToolchain(language)
	default:
		trace.Error("Toolchain \"", name, "\" not supported")
		return nil
	}

	if !verifyCompatibility(toolchain, language) {
		return nil
	}
	return toolchain
}

func defaultToolchain(language Language) Toolchain {
	switch language {
	case JavaLanguage{}:
		return GradleToolchain{}
	case CppLanguage{}:
		return CmakeToolchain{}
	default:
		trace.Error("No supported toolchain for language ", language.Name())
	}
	return nil
}

func verifyCompatibility(toolchain Toolchain, language Language) bool {
	if toolchain == nil || language == nil {
		return false
	}
	if !toolchain.supports(language) {
		trace.Error("Toolchain ", toolchain.name(),
			" does not support language ", language.Name())
		return false
	}
	return true
}

func runBuild(toolchain Toolchain) error {
	wd, _ := os.Getwd()
	buildCommandPath := filepath.Join(wd, toolchain.buildCommandName())
	//trace.Info(buildCommandPath)
	output, err := sh.Command(
		buildCommandPath,
		toolchain.buildCommandArgs()).Output()
	if output != nil {
		trace.Echo(string(output))
	}
	return err
}

// Gradle ========================================================================

func runTests(toolchain Toolchain) error {
	wd, _ := os.Getwd()
	testCommandPath := filepath.Join(wd, toolchain.testCommandName())
	//trace.Info(testCommandPath)
	output, err := sh.Command(
		testCommandPath,
		toolchain.testCommandArgs()).Output()
	if output != nil {
		trace.Echo(string(output))
	}
	return err
}

type GradleToolchain struct {
}

func (toolchain GradleToolchain) name() string {
	return "gradle"
}

func (toolchain GradleToolchain) runBuild() error {
	return runBuild(toolchain)
}

func (toolchain GradleToolchain) runTests() error {
	return runTests(toolchain)
}

func (toolchain GradleToolchain) buildCommandName() string {
	return "gradlew"
}

func (toolchain GradleToolchain) buildCommandArgs() []string {
	return []string{"build", "-x", "test"}
}

func (toolchain GradleToolchain) testCommandName() string {
	return "gradlew"
}

func (toolchain GradleToolchain) testCommandArgs() []string {
	return []string{"test"}
}

// Cmake ========================================================================

func (toolchain GradleToolchain) supports(language Language) bool {
	return language == JavaLanguage{}
}

type CmakeToolchain struct{}

func (toolchain CmakeToolchain) name() string {
	return "cmake"
}

func (toolchain CmakeToolchain) runBuild() error {
	return runBuild(toolchain)
}

func (toolchain CmakeToolchain) runTests() error {
	return runTests(toolchain)
}

func (toolchain CmakeToolchain) buildCommandArgs() []string {
	return []string{"--build", ".", "--config", "Debug"}
}

func (toolchain CmakeToolchain) testCommandArgs() []string {
	return []string{"--output-on-failure", "-C", "Debug"}
}

// Maven ========================================================================

func (toolchain CmakeToolchain) supports(language Language) bool {
	return language == CppLanguage{}
}

type MavenToolchain struct {
}

func (toolchain MavenToolchain) name() string {
	return "maven"
}

func (toolchain MavenToolchain) runBuild() error {
	return runBuild(toolchain)
}

func (toolchain MavenToolchain) runTests() error {
	return runTests(toolchain)
}

func (toolchain MavenToolchain) buildCommandName() string {
	return "mvnw"
}

func (toolchain MavenToolchain) buildCommandArgs() []string {
	return []string{"test-compile"}
}

func (toolchain MavenToolchain) testCommandName() string {
	return "mvnw"
}

func (toolchain MavenToolchain) testCommandArgs() []string {
	return []string{"test"}
}

func (toolchain MavenToolchain) supports(language Language) bool {
	return language == JavaLanguage{}
}
