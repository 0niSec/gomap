package services

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

func GrabBanner(target net.IP, resultsMap map[uint16]string) (map[uint16]string, error) {
	banners := make(map[uint16]string)

	for port, status := range resultsMap {
		if status == "open" {
			address := fmt.Sprintf("%s:%d", target.To4().String(), port)
			conn, err := net.DialTimeout("tcp", address, 5*time.Second)
			if err != nil {
				continue // Skip to the next port
			}
			defer conn.Close()

			err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				continue // Skip to the next port
			}

			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil && err != io.EOF {
				continue // Skip to the next port
			}

			if n > 0 {
				banners[port] = string(buffer[:n])
			}
		}
	}

	return banners, nil
}

//go:embed nmap-services
var nmapServicesData string

// GetServices returns a map of open ports to their corresponding service names
func GetServices(openPorts []uint16) (map[uint16]string, error) {
	services := make(map[uint16]string)

	scanner := bufio.NewScanner(strings.NewReader(nmapServicesData))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		portProto := strings.Split(fields[1], "/")
		if len(portProto) != 2 || portProto[1] != "tcp" {
			continue
		}
		port, err := strconv.ParseUint(portProto[0], 10, 16)
		if err != nil {
			continue
		}
		if contains(openPorts, uint16(port)) {
			services[uint16(port)] = fields[0]
		}
	}

	return services, scanner.Err()
}

// contains returns true if the given port is in the list of ports
func contains(ports []uint16, port uint16) bool {
	for _, p := range ports {
		if p == port {
			return true
		}
	}
	return false
}
