package tracer

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/pion/stun"
)

// check available fields here: https://ip-api.com/#8.8.8.8
type geoResponse struct {
	Country      string `json:"country"`
	Isp          string `json:"isp"`
	QueryStatus  string `json:"status"`
	QueryMessage string `json:"message"`
}

func locateIP(ip string) string {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("[locateIP] failed to fetch location info for IP", ip, err)
		return "N/A"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[locateIP] failed to read resp body:", err)
		return "N/A"
	}

	var geo geoResponse
	err = json.Unmarshal(body, &geo)
	if err != nil {
		fmt.Println("[locateIP] error unmarshaling response:", err)
		return "N/A"
	}

	if geo.QueryStatus == "fail" {
		return geo.QueryMessage
	}

	return fmt.Sprintf("%v (%v)", geo.Isp, geo.Country)
}

// PublicIP returns IPv4, Ipv6 and location of the caller
func PublicIP() (string, string, string) {
	var v4, v6 string
	var v4Err, v6Err error

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		v4, v4Err = stunRequest("udp4")
		wg.Done()
	}()
	go func() {
		v6, v6Err = stunRequest("udp6")
		wg.Done()
	}()

	wg.Wait()

	if v4Err != nil && v4 == "" {
		v4 = "N/A"
	}
	if v6Err != nil && v6 == "" {
		v6 = "N/A"
	}

	return v4, v6, locateIP(v4)
}

// protocol = udp4, udp6
func stunRequest(protocol string) (string, error) {
	var pubIP net.IP
	var stunErr error
	serverAddr, err := net.ResolveUDPAddr(protocol, "stun.l.google.com:19302")
	if err != nil {
		return "", fmt.Errorf("failed to resolve STUN server address: %v", err)
	}

	conn, err := net.DialUDP(protocol, nil, serverAddr)
	if err != nil {
		return "", fmt.Errorf("failed to dial STUN server: %v", err)
	}
	defer conn.Close()

	// Create a new STUN client
	c, err := stun.NewClient(conn)
	if err != nil {
		return "", fmt.Errorf("failed to create STUN client: %v", err)
	}
	defer c.Close()

	// Building binding request with random transaction id.
	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	// Sending request to STUN server, waiting for response message.
	if err := c.Do(message, func(res stun.Event) {
		if res.Error != nil {
			stunErr = res.Error
			return
		}
		// Decoding XOR-MAPPED-ADDRESS attribute from message.
		var xorAddr stun.XORMappedAddress
		if err := xorAddr.GetFrom(res.Message); err != nil {
			stunErr = res.Error
			return
		}
		// could be IPv4 or IPv6 depending on the protocol
		pubIP = xorAddr.IP

	}); err != nil {
		return "", err
	}

	return pubIP.String(), stunErr
}
