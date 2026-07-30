package relay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type InternalAgentAuthenticator struct {
	baseURL string
	secret  string
	client  *http.Client
}

func NewInternalAgentAuthenticator(baseURL, secret string, client *http.Client) (*InternalAgentAuthenticator, error) {
	if strings.TrimSpace(baseURL) == "" || strings.TrimSpace(secret) == "" {
		return nil, fmt.Errorf("internal api url and secret are required")
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &InternalAgentAuthenticator{baseURL: strings.TrimRight(baseURL, "/"), secret: secret, client: client}, nil
}

func (a *InternalAgentAuthenticator) Authenticate(ctx context.Context, token string) (AgentIdentity, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, a.baseURL+"/internal/agents/authenticate", nil)
	if err != nil {
		return AgentIdentity{}, fmt.Errorf("create agent authentication request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Internal-Secret", a.secret)
	response, err := a.client.Do(request)
	if err != nil {
		return AgentIdentity{}, fmt.Errorf("authenticate agent with api: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return AgentIdentity{}, fmt.Errorf("agent authentication returned status %d", response.StatusCode)
	}
	var body struct {
		AgentID        string `json:"agentId"`
		OrganizationID string `json:"organizationId"`
		Limits         struct {
			MaxTunnels     int   `json:"maxTunnels"`
			MaxConnections int   `json:"maxConnections"`
			BandwidthBytes int64 `json:"bandwidthBytes"`
		} `json:"limits"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return AgentIdentity{}, fmt.Errorf("decode agent authentication response: %w", err)
	}
	if body.AgentID == "" || body.OrganizationID == "" {
		return AgentIdentity{}, fmt.Errorf("agent authentication response is incomplete")
	}
	return AgentIdentity{AgentID: body.AgentID, OrganizationID: body.OrganizationID, MaxTunnels: body.Limits.MaxTunnels, MaxConnections: body.Limits.MaxConnections, BandwidthBytes: body.Limits.BandwidthBytes}, nil
}
