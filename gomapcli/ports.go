// This file contains any functions related to ports for the CLI

package gomapcli

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func ParseTarget(target string) (string, error) {
	// Check if the target is an IP address
	// If it is, return it
	ip := net.ParseIP(target)

	if ip != nil {
		return target, nil
	}

	// Check if the target is a domain name
	// If it is, return the IP address
	ips, err := net.LookupIP(target)

	if err != nil {
		return "", fmt.Errorf("error looking up IP address for target '%s': %w", target, err)
	}

	if len(ips) == 0 {
		return "", fmt.Errorf("no IP addresses found for target '%s'", target)
	}

	return ips[0].String(), nil
}

func ParsePorts(ports string) ([]int, error) {
	before, after, found := strings.Cut(ports, "-")
	portRange := []int{}

	if found {
		before, err := strconv.Atoi(before)

		if err != nil {
			return nil, fmt.Errorf("error converting 'before' part of range '%d' to int: %w", before, err)
		}

		after, err := strconv.Atoi(after)

		if err != nil {
			return nil, fmt.Errorf("error converting 'before' part of range '%d' to int: %w", after, err)
		}

		for i := before; i <= after; i++ {
			portRange = append(portRange, i)
		}

		return portRange, nil

	} else {
		singlePort, err := strconv.Atoi(ports)

		if err != nil {
			return nil, fmt.Errorf("error converting single port '%d' to int: %w", singlePort, err)
		}

		return []int{singlePort}, nil
	}
}

func ScanPorts(ports []int, target string, timeout time.Duration) {
	// Get the current time and format it
	// https://pkg.go.dev/time#Time.Format
	startTime := time.Now().Local().Format("2006-01-02 15:04:05")
	openPorts := []int{}
	filteredPorts := []int{}
	closedPorts := []int{}

	fmt.Printf("Starting gomap at %s on target %s\n", startTime, target)

	// Iterate over the ports and scan them
	for _, port := range ports {
		// Construct the address to pass into DialTimeout
		// DialTimeout accepts a string in the format "host:port"
		address := fmt.Sprintf("%s:%d", target, port)
		conn, err := net.DialTimeout("tcp", address, timeout)

		// Check if the error is a timeout error
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			filteredPorts = append(filteredPorts, port)
			continue
		} else if err != nil {
			closedPorts = append(closedPorts, port)
			continue
		} else {
			openPorts = append(openPorts, port)
			conn.Close()
		}

	}

	// Print the closed ports
	if len(closedPorts) > 0 {
		fmt.Printf("Not shown: %d closed tcp ports\n", len(closedPorts))
	}

	if len(filteredPorts) > 0 {
		fmt.Printf("%d filtered ports\n", len(filteredPorts))
	}

	// Print the open ports
	if len(openPorts) == 0 {
		fmt.Println("No open ports found")
		return
	} else {
		fmt.Printf("%-10s %-10s %-10s\n", "PORT", "STATE", "SERVICE")

		for _, port := range openPorts {
			fmt.Printf("%-10s %-10s\n", fmt.Sprintf("%d/tcp", port), "open")
		}
	}

	fmt.Printf("\ngomap completed at %s.\n", time.Now().Local().Format("2006-01-02 15:04:05"))
}
