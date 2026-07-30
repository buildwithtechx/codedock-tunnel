package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CloudflareConfig struct {
	BaseURL    string
	ZoneID     string
	APIToken   string
	HTTPClient *http.Client
}
type Cloudflare struct {
	baseURL string
	zoneID  string
	token   string
	client  *http.Client
}

func NewCloudflare(cfg CloudflareConfig) (*Cloudflare, error) {
	if cfg.ZoneID == "" || cfg.APIToken == "" {
		return nil, fmt.Errorf("cloudflare zone and api token are required")
	}
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.cloudflare.com/client/v4"
	}
	client := cfg.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	return &Cloudflare{baseURL: baseURL, zoneID: cfg.ZoneID, token: cfg.APIToken, client: client}, nil
}

func (c *Cloudflare) UpsertTXT(ctx context.Context, name, value string, ttl int) error {
	if name == "" || value == "" || ttl <= 0 {
		return fmt.Errorf("dns record name, value, and ttl are required")
	}
	return c.request(ctx, http.MethodPost, "/zones/"+c.zoneID+"/dns_records", map[string]any{"type": "TXT", "name": name, "content": value, "ttl": ttl})
}

func (c *Cloudflare) request(ctx context.Context, method, path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode dns request: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create dns request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+c.token)
	request.Header.Set("Content-Type", "application/json")
	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("send dns request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("dns provider returned status %d", response.StatusCode)
	}
	return nil
}
