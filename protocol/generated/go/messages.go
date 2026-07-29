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
	MessageTypeHTTPRequest   MessageType = "http_request"
	MessageTypeHTTPResponse  MessageType = "http_response"
	MessageTypeTCPData       MessageType = "tcp_data"
	MessageTypeTCPClose      MessageType = "tcp_close"
	MessageTypeUDPData       MessageType = "udp_data"
	MessageTypeUDPResponse   MessageType = "udp_response"
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

type HTTPRequest struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body,omitempty"`
}

type HTTPResponse struct {
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body,omitempty"`
	Error      string              `json:"error,omitempty"`
}

type TCPData struct {
	ConnectionID string `json:"connection_id"`
	Data         string `json:"data"`
}

type TCPClose struct {
	ConnectionID string `json:"connection_id"`
	Reason       string `json:"reason,omitempty"`
}

type UDPData struct {
	PacketID      string `json:"packet_id"`
	SourceAddress string `json:"source_address"`
	SourcePort    int    `json:"source_port"`
	Data          string `json:"data"`
}

type UDPResponse struct {
	PacketID      string `json:"packet_id"`
	TargetAddress string `json:"target_address"`
	TargetPort    int    `json:"target_port"`
	Data          string `json:"data"`
}
