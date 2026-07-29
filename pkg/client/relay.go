package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"codedock.run/codedock-tunnel/pkg/protocol"
	"github.com/gorilla/websocket"
)

type RelayConfig struct {
	URL         string
	Token       string
	HTTPHeaders http.Header
}

type RelayConnection struct {
	conn      *websocket.Conn
	TunnelID  string
	PublicURL string
	writeMu   sync.Mutex
}

func OpenRelay(ctx context.Context, cfg RelayConfig, open protocol.OpenTunnel) (*RelayConnection, error) {
	if strings.TrimSpace(cfg.URL) == "" || strings.TrimSpace(cfg.Token) == "" {
		return nil, fmt.Errorf("relay url and token are required")
	}
	if open.LocalPort < 1 || open.LocalPort > 65535 || open.Protocol == "" {
		return nil, fmt.Errorf("valid local port and protocol are required")
	}
	headers := cfg.HTTPHeaders
	if headers == nil {
		headers = make(http.Header)
	}
	headers.Set("Authorization", "Bearer "+cfg.Token)
	connection, _, err := websocket.DefaultDialer.DialContext(ctx, cfg.URL, headers)
	if err != nil {
		return nil, fmt.Errorf("connect relay: %w", err)
	}
	open.Token = cfg.Token
	message, err := protocol.EncodePayload(protocol.MessageTypeOpenTunnel, "", open)
	if err != nil {
		connection.Close()
		return nil, err
	}
	if err := connection.WriteMessage(websocket.TextMessage, message); err != nil {
		connection.Close()
		return nil, fmt.Errorf("send tunnel open: %w", err)
	}
	_, response, err := connection.ReadMessage()
	if err != nil {
		connection.Close()
		return nil, fmt.Errorf("read tunnel acknowledgement: %w", err)
	}
	envelope, err := protocol.Decode(response)
	if err != nil {
		connection.Close()
		return nil, err
	}
	if envelope.Type == protocol.MessageTypeError {
		var failure protocol.ErrorMessage
		_ = protocol.DecodePayload(envelope, &failure)
		connection.Close()
		return nil, fmt.Errorf("relay rejected tunnel: %s", failure.Message)
	}
	if envelope.Type != protocol.MessageTypeOpenTunnelAck {
		connection.Close()
		return nil, fmt.Errorf("unexpected relay response %q", envelope.Type)
	}
	var ack protocol.OpenTunnelAck
	if err := protocol.DecodePayload(envelope, &ack); err != nil {
		connection.Close()
		return nil, err
	}
	return &RelayConnection{conn: connection, TunnelID: ack.TunnelID, PublicURL: ack.PublicURL}, nil
}

func (c *RelayConnection) SendHeartbeat() error {
	if c == nil || c.conn == nil {
		return fmt.Errorf("relay connection is required")
	}
	message, err := protocol.EncodePayload(protocol.MessageTypeHeartbeat, "", protocol.Heartbeat{Timestamp: time.Now().Unix()})
	if err != nil {
		return err
	}
	return c.write(message)
}

func (c *RelayConnection) ServeLocal(ctx context.Context, targetURL string) error {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read relay message: %w", err)
		}
		message, err := protocol.Decode(data)
		if err != nil {
			return err
		}
		if message.Type != protocol.MessageTypeHTTPRequest {
			continue
		}
		var request protocol.HTTPRequest
		if err := protocol.DecodePayload(message, &request); err != nil {
			return err
		}
		response := c.forwardHTTP(ctx, targetURL, request)
		payload, err := protocol.EncodePayload(protocol.MessageTypeHTTPResponse, message.RequestID, response)
		if err != nil {
			return err
		}
		if err := c.write(payload); err != nil {
			return err
		}
	}
}

func (c *RelayConnection) forwardHTTP(ctx context.Context, targetURL string, incoming protocol.HTTPRequest) protocol.HTTPResponse {
	body, err := base64.StdEncoding.DecodeString(incoming.Body)
	if err != nil {
		return protocol.HTTPResponse{Error: "invalid request body"}
	}
	request, err := http.NewRequestWithContext(ctx, incoming.Method, strings.TrimRight(targetURL, "/")+incoming.Path, strings.NewReader(string(body)))
	if err != nil {
		return protocol.HTTPResponse{Error: err.Error()}
	}
	for key, values := range incoming.Headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return protocol.HTTPResponse{Error: err.Error()}
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, 16<<20))
	if err != nil {
		return protocol.HTTPResponse{Error: err.Error()}
	}
	headers := make(map[string][]string, len(response.Header))
	for key, values := range response.Header {
		headers[key] = append([]string(nil), values...)
	}
	return protocol.HTTPResponse{StatusCode: response.StatusCode, Headers: headers, Body: base64.StdEncoding.EncodeToString(data)}
}

func (c *RelayConnection) write(data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *RelayConnection) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
