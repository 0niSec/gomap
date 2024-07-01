// The gomapcli package contains the main logic for the CLI application.
package gomapcli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Runner(c *cli.Context) error {
	// If the ports flag is not set, scan the top 1000 ports by default
	// We only need to parse the target in this case ince the ports are already in the proper format
	if c.String("ports") == "" {
		target, err := ParseTarget(c.String("target"))
		if err != nil {
			return fmt.Errorf("error parsing target: %w", err)
		}

		ScanPorts(Top1000Ports[:], target, c.Duration("timeout"))

		return nil
	}

	// If the ports flag is set, scan the ports specified
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
