package run

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type HelpTestSuite struct {
	suite.Suite
}

func TestHelpTestSuite(t *testing.T) {
	suite.Run(t, new(HelpTestSuite))
}

func (suite *HelpTestSuite) TestNewHelp() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mCmd := NewMockcmd(ctrl)
	originalCommand := command

	command = func(path string, args ...string) cmd {
		assert.Equal(suite.T(), helmBin, path)
		assert.Equal(suite.T(), []string{"help"}, args)
		return mCmd
	}
	defer func() { command = originalCommand }()

	mCmd.EXPECT().
		Stdout(gomock.Any())
	mCmd.EXPECT().
		Stderr(gomock.Any())
	mCmd.EXPECT().
		Run().
		Times(1)

	h := NewHelp(false)
	h.Run()
}
