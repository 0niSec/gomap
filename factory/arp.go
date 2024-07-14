package factory

import (
	"fmt"
	"net"

	"github.com/0niSec/gomap/logger"
	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
)

// ? We may not need this after all
// The initial thought was that we'd need this to get the MAC address of the default gateway
// This is seemingly handled by the network stack since we're operating at the IP Layer

// CreateARPPacket constructs an ARP packet to determine the MAC Address of the default gateway on the network
//
// ARP is used in gomap to determine the MAC Address of the default gateway on the network and the MAC address of local targets
func CreateARPPacket(srcMAC net.HardwareAddr, srcIP, dstIP net.IP) ([]byte, error) {
	// Ethernet Header
	ethHeader := layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	// ARP Header
	arpHeader := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   srcMAC,
		SourceProtAddress: srcIP,
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0}, // Unknown target MAC
		DstProtAddress:    dstIP,
	}

	// Create a buffer to store the packet
	buffer := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	// Serialize the layers into the buffer
	err := gopacket.SerializeLayers(buffer, opts, &ethHeader, &arpHeader)
	if err != nil {
		logger.Error("Failed to serialize layers", "err", err)
		return nil, fmt.Errorf("error serializing layers: %w", err)
	}

	return buffer.Bytes(), nil
}
