// The factory package contains functions for constructing packets.
package factory

import (
	"fmt"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// ConstructICMPPacket constructs an ICMP packet and returns the bytes
// Uses the [net/ipv4] and [net/icmp] packages
func ConstructICMPPacket() ([]byte, error) {
	// Create a new ICMP message by constructing the message struct
	// https://pkg.go.dev/golang.org/x/net/icmp#Message
	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   0,
			Seq:  0,
			Data: []byte("PING"),
		},
	}

	// Marshal the message into bytes
	// https://pkg.go.dev/golang.org/x/net/icmp#Message.Marshal
	messageBytes, err := message.Marshal(nil)

	// ? Is this the best way to do this?
	if err != nil {
		return nil, fmt.Errorf("error marshalling ICMP message: %v", err)
	}

	return messageBytes, nil
}
