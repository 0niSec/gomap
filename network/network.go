package network

import (
	"fmt"
	"net"
	"time"

	"github.com/0niSec/gomap/factory"
	"github.com/0niSec/gomap/logger"
	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcap"
)

// GetValidInterface returns a valid interface used for packet sending and capture
func GetValidInterface() (*net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Error("Failed to get interfaces", "err", err)
		return nil, fmt.Errorf("error getting interfaces: %w", err)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			logger.Error("Failed to get interface addresses %s %w", iface.Name, err)
			return nil, fmt.Errorf("error getting interface addresses: %w", err)
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return &iface, nil // Found a valid IPv4 interface
				}
			}
		}
	}

	logger.Error("No valid interfaces found")
	return nil, fmt.Errorf("no valid interfaces found")
}

// GetMACAddress returns the MAC address of a given target IP address and uses the given interface to send the ARP request
func GetMACAddress(iface *net.Interface, target net.IP) (net.HardwareAddr, error) {
	// Ensure we're using Ipv4
	target = target.To4()
	if target == nil {
		logger.Error("Invalid target IP address", "target", target)
		return nil, fmt.Errorf("invalid target IP address: %s", target)
	}

	// Open a handle to the interface
	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		logger.Error("Failed to open interface", "err", err)
		return nil, fmt.Errorf("error opening interface: %w", err)
	}
	defer handle.Close()

	// Get the source IP address for the interface
	addrs, err := iface.Addrs()
	if err != nil {
		logger.Error("Failed to get interface addresses", "err", err)
		return nil, fmt.Errorf("error getting interface addresses: %w", err)
	}

	// Get the source IP address for the interface
	var srcIP net.IP
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				srcIP = ip4
				break
			}
		}
	}

	if srcIP == nil {
		return nil, fmt.Errorf("no IPv4 address found for interface %s", iface.Name)
	}

	arpRequest, err := factory.CreateARPPacket(iface.HardwareAddr, srcIP, target)
	if err != nil {
		return nil, fmt.Errorf("error creating ARP request: %w", err)
	}

	// Send ARP request
	if err := handle.WritePacketData(arpRequest); err != nil {
		return nil, fmt.Errorf("error sending ARP request: %w", err)
	}

	// Listen for ARP reply
	start := time.Now()
	for time.Since(start) < 5*time.Second {
		data, _, err := handle.ReadPacketData()
		if err != nil {
			continue
		}

		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
		arpLayer := packet.Layer(layers.LayerTypeARP)
		if arpLayer == nil {
			continue
		}

		arp := arpLayer.(*layers.ARP)
		if net.IP(arp.SourceProtAddress).Equal(target) && arp.Operation == layers.ARPReply {
			return net.HardwareAddr(arp.SourceHwAddress), nil
		}
	}

	return nil, fmt.Errorf("no ARP reply received from %s", target)
}
