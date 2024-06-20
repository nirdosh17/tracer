package tracer

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

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

func PublicIP() (string, error) {
	var pubIP net.IP
	var stunErr error

	u, err := stun.ParseURI("stun:stun.l.google.com:19302")
	if err != nil {
		return "", err
	}

	// Creating a "connection" to STUN server.
	c, err := stun.DialURI(u, &stun.DialConfig{})
	if err != nil {
		return "", err
	}

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
		pubIP = xorAddr.IP

	}); err != nil {
		return "", err
	}

	return fmt.Sprintf("%v | %v", pubIP, locateIP(pubIP.String())), stunErr
}
