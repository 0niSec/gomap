package scanner

import (
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/0niSec/gomap/factory"
	"github.com/0niSec/gomap/logger"
	"github.com/0niSec/gomap/network"
	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcap"
)

// SendSYNPacket sends a raw TCP SYN packet with the provided packet data, source IP, and destination IP.
// It creates a raw socket, binds it to the appropriate network interface, and sends the packet using the socket.
// This function is used to initiate a TCP connection by sending a SYN packet.
func SendSYNPacket(packetData []byte, srcIP, dstIP net.IP) error {
	// Get the interface used for sending the packet
	// ? Not sure if this is really needed but it's here just in case (1)
	iface, err := network.GetValidInterface()
	if err != nil {
		return fmt.Errorf("error getting valid interface: %w", err)
	}

	// Create a raw socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return fmt.Errorf("failed to create raw socket: %w", err)
	}
	defer syscall.Close(fd)

	// Bind the socket to the interface
	// ? (2)
	err = syscall.BindToDevice(fd, iface.Name)
	if err != nil {
		logger.Error("Failed to bind raw socket to interface", "err", err)
		return fmt.Errorf("failed to bind raw socket to interface: %w", err)
	}

	// Prepare the sockaddr_in structure
	// We call [net/ipv4/To4()] to convert the IP address to a 4-byte array
	// Calling just the bytes of dstIP results in the Ipv4 address being represented as Ipv6
	addr := syscall.SockaddrInet4{
		Port: 0, // The port is already set in the packet
		Addr: [4]byte{dstIP.To4()[0], dstIP.To4()[1], dstIP.To4()[2], dstIP.To4()[3]},
	}

	// Send the packet
	err = syscall.Sendto(fd, packetData, 0, &addr)
	if err != nil {
		return fmt.Errorf("failed to send packet: %w", err)
	}

	return nil
}

// ReadSYNACKResponse reads the response to a SYN packet sent to the specified destination IP and port, and returns the status of the port (open, closed, or filtered).
// It uses a PCAP handle to capture packets on the appropriate network interface, and applies a BPF filter to only capture relevant TCP packets.
// The function will block until a response is received or the specified timeout is reached, at which point it will return "filtered" if no response was received.
func ReadSYNACKResponse(srcIP net.IP, dstIP net.IP, srcPort, dstPort uint16, timeout time.Duration) (string, error) {
	logger.Debug("Starting ReadSYNACKResponse", "srcIP", srcIP, "srcPort", srcPort, "dstIP", dstIP, "dstPort", dstPort)

	// Find the appropriate interface
	iface, err := network.GetValidInterface()
	if err != nil {
		logger.Error("Failed to find interface", "err", err)
		return "", fmt.Errorf("error finding interface: %w", err)
	}

	// Open the device for capturing
	handle, err := pcap.OpenLive(iface.Name, 65536, true, 10*time.Millisecond)
	if err != nil {
		logger.Error("Failed to open device", "err", err)
		return "", fmt.Errorf("error opening device: %w", err)
	}
	defer handle.Close()

	// Set BPF filter to only capture relevant packets for this specific port
	filter := fmt.Sprintf("tcp and src host %s and src port %d and dst host %s and dst port %d",
		dstIP.String(), dstPort, srcIP.String(), srcPort)
	logger.Debug("Setting BPF filter", "filter", filter)
	err = handle.SetBPFFilter(filter)
	if err != nil {
		logger.Error("Failed to set BPF filter", "err", err)
		return "", fmt.Errorf("error setting BPF filter: %w", err)
	}

	// Create packet source
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.NoCopy = true

	// Wait for a packet or timeout
	for {
		select {
		case packet := <-packetSource.Packets():
			if packet == nil {
				logger.Debug("Received nil packet")
				continue
			}

			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			if tcpLayer == nil {
				logger.Debug("Packet does not contain TCP layer")
				continue
			}
			tcp, _ := tcpLayer.(*layers.TCP)

			logger.Debug("TCP flags", "SYN", tcp.SYN, "ACK", tcp.ACK, "RST", tcp.RST, "srcPort", tcp.SrcPort, "dstPort", tcp.DstPort)

			// Verify that the packet is for the correct ports
			if tcp.SrcPort != layers.TCPPort(dstPort) || tcp.DstPort != layers.TCPPort(srcPort) {
				logger.Debug("Received packet for wrong ports", "expectedSrcPort", dstPort, "actualSrcPort", tcp.SrcPort, "expectedDstPort", srcPort, "actualDstPort", tcp.DstPort)
				continue
			}

			if tcp.SYN && tcp.ACK {
				logger.Debug("Port is open", "dstPort", dstPort)
				return "open", nil
			} else if tcp.RST {
				logger.Debug("Port is closed", "dstPort", dstPort)
				return "closed", nil
			}

			logger.Debug("Unexpected packet flags", "dstPort", dstPort)
			return "filtered", nil

		case <-time.After(timeout):
			logger.Debug("Timeout reached", "dstPort", dstPort)
			return "filtered", nil
		}
	}
}

// Scan performs a SYN scan on the given source and destination IP addresses and ports.
// It sends a SYN packet to each destination port and waits for a response.
// The function returns a map of port statuses, where the key is the port number and the value is the status ("open", "closed", "filtered", or "error").
// The scan will timeout after the specified duration.
func Scan(srcIP, dstIP net.IP, ports []uint16, timeout time.Duration) map[uint16]string {
	srcPort, err := factory.GenerateRandomPort()
	if err != nil {
		logger.Error("Failed to generate random port", "err", err)
		return nil
	}

	results := make(map[uint16]string)
	resultChan := make(chan struct {
		port   uint16
		status string
	})

	for _, dstPort := range ports {
		go func(dstPort uint16) {
			logger.Debug("Starting goroutine", "dstPort", dstPort)
			// Create the SYN Packet
			packetData, _, _, err := factory.CreateSYNPacket(srcIP, dstIP, srcPort, dstPort)
			if err != nil {
				logger.Error("Failed to create SYN packet", "err", err)
				resultChan <- struct {
					port   uint16
					status string
				}{dstPort, "error"}
				return
			}

			// Send the SYN Packet
			err = SendSYNPacket(packetData, srcIP, dstIP)
			if err != nil {
				logger.Error("Failed to send SYN packet", "err", err)
				resultChan <- struct {
					port   uint16
					status string
				}{dstPort, "error"}
				return
			}

			status, err := ReadSYNACKResponse(srcIP, dstIP, srcPort, dstPort, timeout)
			if err != nil {
				logger.Error("Failed to read SYN/ACK response", "err", err)
				resultChan <- struct {
					port   uint16
					status string
				}{dstPort, "error"}
				return
			}

			resultChan <- struct {
				port   uint16
				status string
			}{dstPort, status}
		}(dstPort)
	}

	for range ports {
		result, ok := <-resultChan
		if !ok {
			logger.Error("Failed to receive result from resultChan")
			break
		}
		results[result.port] = result.status
	}

	return results
}
