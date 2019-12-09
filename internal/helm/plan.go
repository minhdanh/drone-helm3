package helm

import (
	"errors"
	"github.com/pelotech/drone-helm3/internal/run"
	"io"
	"os"
)

// A Step is one step in the plan.
// TODO: If no further interface methods materialize, Step should be called Runner instead
type Step interface {
	Run() error
}

// A Plan is a series of steps to perform.
type Plan struct {
	steps []Step
}

// NewPlan makes a plan for running a helm operation.
func NewPlan(cfg Config) (*Plan, error) {
	gCfg := run.GlobalConfig{
		Debug:          cfg.Debug,
		KubeConfig:     cfg.KubeConfig,
		Values:         cfg.Values,
		StringValues:   cfg.StringValues,
		ValuesFiles:    cfg.ValuesFiles,
		Namespace:      cfg.Namespace,
		Token:          cfg.Token,
		SkipTLSVerify:  cfg.SkipTLSVerify,
		Certificate:    cfg.Certificate,
		APIServer:      cfg.APIServer,
		ServiceAccount: cfg.ServiceAccount,
		Stdout:         os.Stdout,
		Stderr:         os.Stderr,
	}

	p := Plan{}
	switch cfg.Command {
	case "upgrade":
		steps, err := upgrade(cfg, gCfg)
		if err != nil {
			return nil, err
		}
		p.steps = steps
	case "delete":
		return nil, errors.New("not implemented")
	case "lint":
		return nil, errors.New("not implemented")
	case "help":
		return nil, errors.New("not implemented")
	default:
		switch cfg.DroneEvent {
		case "push", "tag", "deployment", "pull_request", "promote", "rollback":
			steps, err := upgrade(cfg)
			if err != nil {
				return nil, err
			}
			p.steps = steps
		default:
			return nil, errors.New("not implemented")
		}
	}

	return &p, nil
}

// Execute runs each step in the plan, aborting and reporting on error
func (p *Plan) Execute() error {
	for _, step := range p.steps {
		if err := step.Run(); err != nil {
			// TODO: fmt.Errorf
			return err
		}
	}

	return nil
}

func upgrade(cfg Config, gCfg run.GlobalConfig) ([]Step, error) {
	uCfg := run.UpgradeConfig{
		Chart:        cfg.Chart,
		Release:      cfg.Release,
		ChartVersion: cfg.ChartVersion,
		Wait:         cfg.Wait,
		ReuseValues:  cfg.ReuseValues,
		Timeout:      cfg.Timeout,
		Force:        cfg.Force,
	}
	steps := make([]Step, 0)
	steps = append(steps, run.NewUpgrade(uCfg, gCfg))

	return steps, nil
}
