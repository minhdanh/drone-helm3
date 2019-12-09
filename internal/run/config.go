package run

import (
	"io"
)

// GlobalConfig is global in the sense of "applies to all helm commands," not "everywhere in the code"
type GlobalConfig struct {
	Debug          bool
	KubeConfig     string
	Values         string
	StringValues   string
	ValuesFiles    []string
	Namespace      string
	Token          string
	SkipTLSVerify  bool
	Certificate    string
	APIServer      string
	ServiceAccount string
	Stdout         io.Writer
	Stderr         io.Writer
}
