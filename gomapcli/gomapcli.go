// The gomapcli package contains the main logic for the CLI application.
package gomapcli

import (
	"fmt"

	"github.com/0niSec/gomap/network"
	"github.com/0niSec/gomap/scanner"
	"github.com/urfave/cli/v2"
)

func Runner(c *cli.Context) error {

	iface, err := network.GetValidInterface()
	if err != nil {
		return fmt.Errorf("error getting valid interface: %w", err)
	}

	// Get IP Address for interface
	srcIP, err := network.GetInterfaceIPAddress(iface)
	if err != nil {
		return fmt.Errorf("error getting interface IP address: %w", err)
	}

	target, err := ParseTarget(c.String("target"))
	if err != nil {
		return fmt.Errorf("error parsing target: %w", err)
	}

	// If the ports flag is not set, scan the top 1000 ports by default
	// We only need to parse the target in this case since the ports are already in the proper format
	if c.String("ports") == "" {
		results := scanner.Scan(srcIP, target, Top1000Ports, c.Duration("timeout"))
		for _, result := range results {
			fmt.Println(result)
		}

		return nil
	}

	// If the ports flag is set, scan the ports specified
	ports, err := ParsePorts(c.String("ports"))
	if err != nil {
		return fmt.Errorf("error parsing ports: %w", err)
	}

	results := scanner.Scan(srcIP, target, ports, c.Duration("timeout"))
	for _, result := range results {
		fmt.Println(result)
	}

	return nil
}
