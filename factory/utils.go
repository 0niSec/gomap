package factory

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	minPort uint16 = 49152 // Start of the dynamic/private port range
	maxPort uint16 = 65535 // End of the port range
)

// Create a new random number generator
var (
	rng       = rand.New(rand.NewSource(time.Now().UnixNano()))
	usedPorts = make(map[uint16]bool)
	portMutex sync.Mutex
)

// GenerateRandomPort generates a random, unused TCP port number.
//
// This function will attempt to find an unused port by generating random port
// numbers within the dynamic/private port range (49152-65535) and checking if
// the port is available. It will try up to 100 times before returning an error.
//
// The returned port number is marked as used in the usedPorts map to prevent
// it from being returned by subsequent calls to this function. When you're done
// using the port, call ReleasePort to mark it as available again.
//
// Returns:
//   - The generated port number as a uint16.
//   - An error if no unused port could be found after 100 attempts.
func GenerateRandomPort() (uint16, error) {
	portMutex.Lock()
	defer portMutex.Unlock()

	for tries := 0; tries < 100; tries++ {
		// Use rng.Intn instead of rand.Intn
		port := uint16(rng.Intn(int(maxPort-minPort+1))) + minPort

		if usedPorts[port] {
			continue
		}

		addr := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", addr)

		if err != nil {
			usedPorts[port] = true
			continue
		}

		listener.Close()
		usedPorts[port] = true
		return port, nil
	}

	return 0, fmt.Errorf("unable to find an unused port after 100 attempts")
}

// ReleasePort marks a port as no longer in use, allowing it to be returned
// by future calls to GenerateRandomPort.
//
// This function should be called when you're done using a port that was
// allocated by GenerateRandomPort.
//
// Usage:
//
//	ReleasePort(port)
//
// Parameters:
//   - port: The uint16 port number to be released.
func ReleasePort(port uint16) {
	portMutex.Lock()
	defer portMutex.Unlock()
	delete(usedPorts, port)
}
