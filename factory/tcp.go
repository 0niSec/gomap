package factory

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/0niSec/gomap/logger"
	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
)

type PortStatus int

// CreateSYNPacket creates a TCP SYN packet with the specified source and destination IP and port.
// It returns the serialized packet bytes, the IPv4 layer, and the TCP layer.
// If there is an error generating the packet, it returns an error.
func CreateSYNPacket(srcIP, dstIP net.IP, srcPort, dstPort uint16) ([]byte, *layers.IPv4, *layers.TCP, error) {
	// Create IP Layer
	ipLayer := &layers.IPv4{
		Version: 4,
		TTL:     64,
		IHL:     5,
		// SrcIP and DstIP are left as net.IP types and not converted into IPv4 here
		// They will be converted into IPv4 when the socket connection is made
		// Gopacket probably does this for us but I'm not sure
		SrcIP:    srcIP, // net.IP
		DstIP:    dstIP, // net.IP
		Protocol: layers.IPProtocolTCP,
	}

	// Create TCP Layer
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(srcPort),
		DstPort: layers.TCPPort(dstPort),
		Ack:     0,
		SYN:     true,
		ACK:     false,
		PSH:     false,
		RST:     false,
		URG:     false,
		ECE:     false,
		CWR:     false,
		NS:      false,
		FIN:     false,
		Urgent:  0,
		Window:  65535,
		Seq:     rand.Uint32(),
		Options: []layers.TCPOption{
			{
				OptionType:   layers.TCPOptionKindMSS,
				OptionLength: 4,
				OptionData:   []byte{0x05, 0xb4},
			},
			{
				OptionType:   layers.TCPOptionKindSACKPermitted,
				OptionLength: 2,
			},
			{
				OptionType:   layers.TCPOptionKindTimestamps,
				OptionLength: 10,
				OptionData:   generateTimestampOption(),
			},
			{
				OptionType:   layers.TCPOptionKindNop,
				OptionLength: 1,
			},
			{
				OptionType:   layers.TCPOptionKindWindowScale,
				OptionLength: 3,
				OptionData:   []byte{7},
			},
		},
	}

	// ! DEBUG
	// logger.Debug("Created SYN packet", "packet", tcpLayer)

	// Set TCP Checksum
	err := tcpLayer.SetNetworkLayerForChecksum(ipLayer)
	if err != nil {
		logger.Error("Failed to set network layer for TCP checksum", "err", err)
		return nil, nil, nil, fmt.Errorf("error setting network layer for TCP checksum: %w", err)
	}

	// Serialize the layers into the buffer
	buffer := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	err = gopacket.SerializeLayers(buffer, opts, ipLayer, tcpLayer)
	if err != nil {
		logger.Error("Failed to serialize layers while creating SYN Packet", "err", err)
		return nil, nil, nil, fmt.Errorf("error serializing layers while creating SYN packet: %w", err)
	}

	// ! DEBUG
	// logger.Debug("Packet serialized", "packet", buffer.Bytes())

	return buffer.Bytes(), ipLayer, tcpLayer, nil
}

// generateTimestampOption generates the TCP timestamp option data, which includes the current timestamp and an echo reply of 0 for the initial SYN packet.
func generateTimestampOption() []byte {
	tsVal := make([]byte, 4)
	tsEcr := make([]byte, 4)

	// Current Timestamp
	binary.BigEndian.PutUint32(tsVal, uint32(time.Now().Unix()))

	// Echo reply is 0 for initial SYN
	binary.BigEndian.PutUint32(tsEcr, 0)

	return append(tsVal, tsEcr...)
}
