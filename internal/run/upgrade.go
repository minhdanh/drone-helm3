package run

import (
	"fmt"
)

// UpgradeConfig has configuration specific to the `helm upgrade` command.
type UpgradeConfig struct {
	Chart   string
	Release string

	ChartVersion string
	Wait         bool
	ReuseValues  bool
	Timeout      string
	Force        bool
}

// Upgrade is an execution step that calls `helm upgrade` when it runs.
type Upgrade struct {
	cfg  UpgradeConfig
	gCfg GlobalConfig
	cmd  cmd
}

// Run executes the `helm upgrade` command.
func (u *Upgrade) Run() error {
	return u.cmd.Run()
}

// NewUpgrade creates a new Upgrade.
func NewUpgrade(cfg UpgradeConfig, gCfg GlobalConfig) *Upgrade {
	u := Upgrade{
		cfg:  cfg,
		gCfg: gCfg,
	}
	args := []string{"upgrade", "--install", cfg.Release, cfg.Chart}
	// TODO: bail if chart/release isn't present? Or just let helm handle that?

	if gCfg.Debug {
		args = append([]string{"--debug"}, args...)
	}

	u.cmd = command(helmBin, args...)
	u.cmd.Stdout(gCfg.Stdout)
	u.cmd.Stderr(gCfg.Stderr)

	if gCfg.Debug {
		fmt.Fprintf(gCfg.Stderr, "Generated command: '%s'\n", u.cmd.String())
	}

	return &u
}
