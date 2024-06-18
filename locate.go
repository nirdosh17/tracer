package tracer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
