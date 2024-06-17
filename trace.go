package tracer

import (
	"fmt"
	"log"
	"net"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	UDPStartPort = 33434
)

type Tracer struct {
	Config *TracerConfig
}

type Hop struct {
	Addr     string
	Location string
	// current ttl(hops) of the packet
	TTL int
	// total time taken for this hop
	ElapsedTime time.Duration
}

func NewTracer(c *TracerConfig) *Tracer {
	return &Tracer{Config: c}
}

type NetworkTrace struct {
	RoundTripTime time.Duration
	NetworkHops   []Hop
}

// Run sends packets to the specified host in loop recording each network hop until it reaches the destination or max hops is reached.
// It also collects traces in the given channel.
//
// e.g. host = example.com
func (t Tracer) Run(host string, traces chan Hop) (NetworkTrace, error) {
	nTrace := NetworkTrace{}
	nHops := []Hop{}

	// resolve host(e.g. example.com) into an IP
	destIP, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return nTrace, fmt.Errorf("unable to resolve host %s", host)
	}

	roundTripStart := time.Now()
	for ttl := 1; ttl <= t.Config.MaxHops; ttl++ {
		// using different UDP port each time
		port := UDPStartPort + ttl
		addr := fmt.Sprintf("%s:%d", destIP, port)

		now := time.Now()
		err = t.sendUDPPacket(addr, ttl)
		if err != nil {
			fmt.Printf("Error sending UDP packet: %s\n", err)
			continue
		}

		recv, err := t.listenICMPMessages()
		if err != nil {
			fmt.Printf("Error listening for ICMP replies: %s\n", err)
			continue
		}
		elapsedTime := time.Since(now)

		// TODO: find out why this is nil
		if recv.packetAddr == nil {
			continue
		}

		packetAddr := recv.packetAddr.String()
		if recv.destIP == destIP.IP.String() {
			loc := locateIP(packetAddr)
			hop := Hop{
				TTL:         ttl,
				Addr:        packetAddr,
				Location:    fmt.Sprintf("%v (%v)", loc.Country, loc.Org),
				ElapsedTime: elapsedTime,
			}

			// push to channel for live updates
			traces <- hop

			nHops = append(nHops, hop)
		}

		if packetAddr == recv.destIP {
			break
		}
	}
	nTrace.RoundTripTime = time.Since(roundTripStart)
	nTrace.NetworkHops = nHops

	close(traces)

	return nTrace, nil
}

// sendUDPPacket sends UDP datagrams with a specified TTL.
// After setting up an ICMP listener, we use this method to send UDP datagrams wrapped in IP packets.
func (t Tracer) sendUDPPacket(addr string, ttl int) error {
	// Resolve the UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Set the TTL
	rawConn, err := conn.SyscallConn()
	if err != nil {
		return err
	}

	err = rawConn.Control(func(fd uintptr) {
		if udpAddr.IP.To4() != nil {
			err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TTL, ttl)
		} else {
			err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, syscall.IPV6_UNICAST_HOPS, ttl)
		}
	})
	if err != nil {
		return err
	}

	// Sending a UDP packet with dummy payload. Why?
	// - actual content in payload sent in the UDP packet does not affect traceroute operation
	// - primary purpose of payload is to generate ICMP responses at each hop
	_, err = conn.Write([]byte("HI"))
	return err
}

// listenICMPMessages listens all ICMP messages incoming in the machine.
// Filter outs unknown messages using caller IP.
func (t Tracer) listenICMPMessages() (icmpResp, error) {
	defaultResp := icmpResp{}
	c, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return icmpResp{}, fmt.Errorf("failed to listen icmp %v", err)
	}
	defer c.Close()

	c.SetReadDeadline(time.Now().Add(time.Duration(t.Config.TimeoutSeconds * int(time.Second))))
	buffer := make([]byte, 1500)

	for {
		receivedBytesLen, receivedFrom, err := c.ReadFrom(buffer)
		if err != nil && receivedBytesLen == 0 {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// after finishing reading from connection, it will timeout in the next loop when there is nothing to read
				// so not need to log error message
				// fmt.Println("Request timed out.", err)
				return defaultResp, nil
			}
			return defaultResp, err
		}

		icmpMsg, err := t.parseICMP(buffer, receivedBytesLen)
		if err != nil {
			return defaultResp, nil
		}

		switch icmpMsg.icmpType {
		case ipv4.ICMPTypeTimeExceeded, ipv4.ICMPTypeDestinationUnreachable:
			icmpMsg.packetAddr = receivedFrom
			return icmpMsg, nil
		default:
			log.Printf("got %+v; want echo reply", icmpMsg.icmpType)
			return defaultResp, nil
		}
	}

}

type icmpResp struct {
	code int
	// the last host where test UDP datagram reached with given TTL or the ICMP sender
	packetAddr  net.Addr
	icmpType    icmp.Type
	requesterIP string
	// we are only interested in ICMP packets which were send to this IP
	// but they might not have reached to this destination due to small TTL
	destIP string
}

func (t Tracer) parseICMP(buffer []byte, length int) (icmpResp, error) {
	var msg icmpResp
	msg.destIP = net.IP(buffer[24:28]).String()
	msg.requesterIP = net.IP(buffer[20:24]).String()

	im, err := icmp.ParseMessage(1, buffer[:length])
	if err != nil {
		return msg, fmt.Errorf("failed to parse icmp msg %v", err)
	}
	msg.icmpType = im.Type
	msg.code = im.Code

	return msg, nil
}
