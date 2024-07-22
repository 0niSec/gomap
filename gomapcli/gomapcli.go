// The gomapcli package contains the main logic for the CLI application.
package gomapcli

import (
	"fmt"
	"time"

	"github.com/0niSec/gomap/network"
	"github.com/0niSec/gomap/scanner"
	"github.com/0niSec/gomap/services"
	"github.com/urfave/cli/v2"
)

func Runner(c *cli.Context) error {
	// Calculate start time
	startTime := time.Now()

	// Get the interface
	iface, err := network.GetValidInterface()
	if err != nil {
		return fmt.Errorf("error getting valid interface: %w", err)
	}

	// Get the source IP for the interface
	srcIP, err := network.GetInterfaceIPAddress(iface)
	if err != nil {
		return fmt.Errorf("error getting interface IP address: %w", err)
	}

	// Parse the target (IP or domain)
	target, err := ParseTarget(c.String("target"))
	if err != nil {
		return fmt.Errorf("error parsing target: %w", err)
	}

	// Parse the ports depending on the -p flag
	ports := Top1000Ports
	if c.String("ports") != "" {
		ports, err = ParsePorts(c.String("ports"))
		if err != nil {
			return fmt.Errorf("error parsing ports: %w", err)
		}
	}

	// Scan the ports
	results, err := scanner.Scan(srcIP, target, ports, c.Duration("timeout"))
	if err != nil {
		return fmt.Errorf("error scanning ports: %w", err)
	}

	// Get the services
	services, err := services.GetServices(ports)
	if err != nil {
		return fmt.Errorf("error loading nmap services: %w", err)
	}

	scanner.PrettyPrintScanResults(results, services)

	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()
	fmt.Printf("Scan completed in %.2f seconds\n", duration)

	return nil
}
