// The gomapcli package contains the main logic for the CLI application.
package gomapcli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Runner(c *cli.Context) error {
	// Parse the ports from the command line
	ports, err := ParsePorts(c.String("ports"))
	if err != nil {
		return fmt.Errorf("error parsing ports: %w", err)
	}

	// Parse the target from the command line
	target, err := ParseTarget(c.String("target"))
	if err != nil {
		return fmt.Errorf("error parsing target: %w", err)
	}

	ScanPorts(ports, target, c.Duration("timeout"))

	return nil
}
