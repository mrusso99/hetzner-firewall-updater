package firewall

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func getIp(ctx context.Context) (net.IP, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://icanhazip.com", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching public IP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	ip := net.ParseIP(string(bytes.TrimSpace(body)))
	if ip == nil {
		return nil, fmt.Errorf("invalid IP format: %s", body)
	}

	log.Println("Public IP:", ip.String())
	return ip, nil
}
