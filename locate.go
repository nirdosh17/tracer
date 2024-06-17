package tracer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// geoResponse represents the response from the geolocation API
type geoResponse struct {
	IP      string `json:"query"`
	Country string `json:"country"`
	Region  string `json:"regionName"`
	City    string `json:"city"`
	Org     string `json:"org"`
}

func locateIP(ip string) geoResponse {
	defaultResp := geoResponse{}
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("[locateIP] failed to fetch location info for IP", ip, err)
		return defaultResp
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[locateIP] failed to read resp body:", err)
		return defaultResp
	}

	var geo geoResponse
	err = json.Unmarshal(body, &geo)
	if err != nil {
		fmt.Println("[locateIP] error unmarshaling response:", err)
	}
	return geo
}
