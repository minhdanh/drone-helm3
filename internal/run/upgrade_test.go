package run

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type UpgradeTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockCmd         *Mockcmd
	originalCommand func(string, ...string) cmd
}

func (suite *UpgradeTestSuite) BeforeTest(_, _ string) {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockCmd = NewMockcmd(suite.ctrl)

	suite.originalCommand = command
	command = func(path string, args ...string) cmd { return suite.mockCmd }
}

func (suite *UpgradeTestSuite) AfterTest(_, _ string) {
	command = suite.originalCommand
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (suite *UpgradeTestSuite) TestNewUpgrade() {
	defer suite.ctrl.Finish()

	cfg := UpgradeConfig{
		Chart:   "at40",
		Release: "jonas_brothers_only_human",
	}

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"upgrade", "--install", "jonas_brothers_only_human", "at40"}, args)

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().
		Stdout(gomock.Any())
	suite.mockCmd.EXPECT().
		Stderr(gomock.Any())
	suite.mockCmd.EXPECT().
		Run().
		Times(1)

	u := NewUpgrade(cfg, GlobalConfig{})
	u.Run()
}

func (suite *UpgradeTestSuite) TestNewUpgradeDebugFlag() {
	cfg := UpgradeConfig{
		Chart:   "at40",
		Release: "lewis_capaldi_someone_you_loved",
	}

	stdout := strings.Builder{}
	stderr := strings.Builder{}
	gCfg := GlobalConfig{
		Debug:  true,
		Stdout: &stdout,
		Stderr: &stderr,
	}

	command = func(path string, args ...string) cmd {
		suite.mockCmd.EXPECT().
			String().
			Return(fmt.Sprintf("%s %s", path, strings.Join(args, " ")))

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().
		Stdout(&stdout)
	suite.mockCmd.EXPECT().
		Stderr(&stderr)

	u := NewUpgrade(cfg, gCfg)

	_ = u

	want := fmt.Sprintf("Generated command: '%s --debug upgrade --install lewis_capaldi_someone_you_loved at40'\n", helmBin)
	suite.Equal(want, stderr.String())
	suite.Equal("", stdout.String())
}
