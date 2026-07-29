package protocol

import "encoding/json"

const Version = 1

type MessageType string

const (
	MessageTypeOpenTunnel    MessageType = "open_tunnel"
	MessageTypeOpenTunnelAck MessageType = "open_tunnel_ack"
	MessageTypeCloseTunnel   MessageType = "close_tunnel"
	MessageTypeData          MessageType = "data"
	MessageTypeHeartbeat     MessageType = "heartbeat"
	MessageTypeError         MessageType = "error"
)

type Envelope struct {
	Version   int             `json:"version"`
	Type      MessageType     `json:"type"`
	RequestID string          `json:"request_id,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

type OpenTunnel struct {
	Token        string `json:"token"`
	LocalPort    int    `json:"local_port"`
	Subdomain    string `json:"subdomain,omitempty"`
	Protocol     string `json:"protocol"`
	CustomDomain string `json:"custom_domain,omitempty"`
}

type OpenTunnelAck struct {
	TunnelID  string `json:"tunnel_id"`
	PublicURL string `json:"public_url"`
}

type CloseTunnel struct {
	TunnelID string `json:"tunnel_id"`
	Reason   string `json:"reason,omitempty"`
}

type Data struct {
	TunnelID string `json:"tunnel_id"`
	StreamID string `json:"stream_id"`
	Data     []byte `json:"data"`
}

type Heartbeat struct {
	Timestamp int64 `json:"timestamp"`
}

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
