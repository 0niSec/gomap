// The factory package contains functions for constructing packets.
package factory

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// CreateICMPPacket constructs an ICMP packet and returns the bytes
// Uses the [net/ipv4] and [net/icmp] packages
func CreateICMPPacket() ([]byte, error) {
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

// SendICMPRequest sends an ICMP request to the provided target with given timeout duration
// and returns true if it receives a reply
func SendICMPRequest(target net.IP) (bool, error) {
	startTime := time.Now()
	var endTime time.Time

	// Construct the ICMP message
	// ? Don't we want the caller ConstructICMPPacket to handle any errors that go wrong? Why do we need an error here?
	messageBytes, err := CreateICMPPacket()
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

	// Set the read deadline to the specified timeout
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Write the ICMP message to the connection
	// https://pkg.go.dev/golang.org/x/net/icmp#PacketConn.WriteTo
	if _, err := conn.WriteTo(messageBytes, &net.IPAddr{IP: target.To4()}); err != nil {
		return false, fmt.Errorf("error writing ICMP message: %v", err)
	}

	// Prepare to receive the reply
	readBuffer := make([]byte, 1500)
	n, _, err := conn.ReadFrom(readBuffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Printf("Host seems down. It did not respond to our ping after 5 seconds\n\n")
			return false, nil // Failed, timed out
		}
		fmt.Printf("error reading ICMP response: %v\n\n", err)
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
		fmt.Printf("Host is alive (%3.4fs)\n\n", endTime.Sub(startTime).Seconds())
		return true, nil // Success, received echo reply
	}

	return false, nil // Failure, did not receive echo reply
}
