package toolchain

import (
	"errors"
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/mengdaming/tcr/language"
	"github.com/mengdaming/tcr/report"
	"os"
	"path/filepath"
)

type Toolchain interface {
	Name() string
	RunBuild() error
	RunTests() error
	buildCommandName() string
	buildCommandArgs() []string
	testCommandName() string
	testCommandArgs() []string
	supports(lang language.Language) bool
}

func New(name string, lang language.Language) (Toolchain, error) {
	var toolchain Toolchain = nil
	var err error

	switch name {
	case GradleToolchain{}.Name():
		toolchain = GradleToolchain{}
	case MavenToolchain{}.Name():
		toolchain = MavenToolchain{}
	case CmakeToolchain{}.Name():
		toolchain = CmakeToolchain{}
	case "":
		toolchain, err = defaultToolchain(lang)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprint("toolchain \"", name, "\" not supported"))
	}

	comp, err := verifyCompatibility(toolchain, lang)
	if !comp || err != nil {
		return nil, err
	}
	return toolchain, nil
}

func defaultToolchain(lang language.Language) (Toolchain, error) {
	switch lang {
	case language.Java{}:
		return GradleToolchain{}, nil
	case language.Cpp{}:
		return CmakeToolchain{}, nil
	default:
		return nil, errors.New(fmt.Sprint("no supported toolchain for ", lang.Name(), " language"))
	}
}

func verifyCompatibility(toolchain Toolchain, lang language.Language) (bool, error) {
	if toolchain == nil || lang == nil {
		return false, errors.New("toolchain and/or language is unknown")
	}
	if !toolchain.supports(lang) {
		return false, errors.New(fmt.Sprintf(
			"%v toolchain does not support %v language",
			toolchain.Name(), lang.Name(),
		))
	}
	return true, nil
}

func runBuild(toolchain Toolchain) error {
	wd, _ := os.Getwd()
	buildCommandPath := filepath.Join(wd, toolchain.buildCommandName())
	output, err := sh.Command(
		buildCommandPath,
		toolchain.buildCommandArgs()).Output()
	if output != nil {
		report.PostText(string(output))
	}
	return err
}

func runTests(tchn Toolchain) error {
	wd, _ := os.Getwd()
	testCommandPath := filepath.Join(wd, tchn.testCommandName())
	output, err := sh.Command(
		testCommandPath,
		tchn.testCommandArgs()).Output()
	if output != nil {
		report.PostText(string(output))
	}
	return err
}