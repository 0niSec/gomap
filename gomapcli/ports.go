package gomapcli

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/0niSec/gomap/factory"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// SendICMPRequest sends an ICMP request to the provided target with given timeout duration
// and returns true if it receives a reply
func SendICMPRequest(target string, timeout time.Duration) (bool, error) {
	startTime := time.Now()
	var endTime time.Time

	// Construct the ICMP message
	// ? Don't we want the caller ConstructICMPPacket to handle any errors that go wrong? Why do we need an error here?
	messageBytes, err := factory.ConstructICMPPacket()
	if err != nil {
		return false, fmt.Errorf("error constructing ICMP message: %v", err)
	}

	// Listen for incoming ICMP packets addressed to `address`
	// The address used is an empty string, which listens on localhost
	// https://pkg.go.dev/golang.org/x/net/icmp#PacketConn
	conn, err := icmp.ListenPacket("ip4:icmp", "")
	if err != nil {
		fmt.Printf("error creating ICMP connection: %v\n", err)
		return false, err
	}
	defer conn.Close()

	// Write the ICMP message to the connection
	// https://pkg.go.dev/golang.org/x/net/icmp#PacketConn.WriteTo
	if _, err := conn.WriteTo(messageBytes, &net.IPAddr{IP: net.ParseIP(target)}); err != nil {
		return false, fmt.Errorf("error writing ICMP message: %v", err)
	}

	// Prepare to receive the reply
	readBuffer := make([]byte, 1500)
	n, _, err := conn.ReadFrom(readBuffer)
	if err != nil {
		fmt.Printf("error reading ICMP response: %v\n", err)
		return false, err
	}

	// Parse the ICMP message
	// https://pkg.go.dev/golang.org/x/net/icmp#ParseMessage
	parsedMessage, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), readBuffer[:n])
	if err != nil {
		fmt.Printf("error parsing ICMP message: %v\n", err)
		return false, err
	}

	// Check if the reply is an ICMP echo reply
	if parsedMessage.Type == ipv4.ICMPTypeEchoReply {
		endTime = time.Now()
		fmt.Printf("Host is alive (%3.4fs)\n", endTime.Sub(startTime).Seconds())
		return true, nil // Success, received echo reply
	}

	fmt.Println("Host did not reply to ICMP")
	return false, nil // Failure, did not receive echo reply
}

// ParseTarget parses the target string and returns the IP address in string format
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

// ParsePorts parses the ports string and returns a slice of integers
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

// ScanPorts scans the provided ports on the target and prints the results
func ScanPorts(ports []int, target string, timeout time.Duration) {
	// Get the current time and format it
	// https://pkg.go.dev/time#Time.Format
	startTime := time.Now().Local().Format("2006-01-02 15:04:05")
	openPorts := []int{}
	filteredPorts := []int{}
	closedPorts := []int{}

	fmt.Printf("Starting gomap at %s on target %s\n", startTime, target)

	SendICMPRequest(target, timeout)

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

	// TODO: The way we determine the port is filtered at the moment is technically incorrect
	// I'd need to build a more robust way to determine if a port is filtered

	// Print the closed ports
	if len(closedPorts) > 0 {
		fmt.Printf("Not shown: %d closed tcp ports\n", len(closedPorts))
	}

	// Print the filtered ports
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
