package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
	"strings"
)

type Config struct {
	// Configuration for drone-helm itself
	Command subcommand `envconfig:"HELM_COMMAND"` // Helm command to run
}

type subcommand string

// subcommand.Decode checks the given value against the list of known commands and generates a helpful error if the command is unknown.
func (cmd *subcommand) Decode(value string) error {
	known := []string{"upgrade", "delete", "lint", "help"}
	for _, c := range known {
		if value == c {
			*cmd = subcommand(value)
			return nil
		}
	}

	if value == "" {
		return nil
	}
	known[len(known)-1] = fmt.Sprintf("or %s", known[len(known)-1])
	return fmt.Errorf("Unknown command '%s'. If specified, command must be %s.",
		value, strings.Join(known, ", "))
}

func main() {
	var c Config

	if err := envconfig.Process("plugin", &c); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	fmt.Printf("config.Command is '%s'\n", c.Command)
}
