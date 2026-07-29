package relay

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"codedock.run/codedock-tunnel/internal/engine"
	"codedock.run/codedock-tunnel/pkg/protocol"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AgentIdentity struct {
	AgentID        string
	OrganizationID string
}

type AgentAuthenticator interface {
	Authenticate(context.Context, string) (AgentIdentity, error)
}

type connectionState struct {
	negotiated    bool
	authenticated bool
}

type Handler struct {
	authenticator AgentAuthenticator
	sessions      *engine.SessionRegistry
	router        *engine.RequestRouter
	tcp           *TCPManager
	udp           *UDPManager
	maxSessions   int
	maxTunnels    int
	maxBandwidth  int64
	heartbeat     time.Duration
	readTimeout   time.Duration
	maxFrameBytes int64
	drainTimeout  time.Duration
	logger        *slog.Logger
	metrics       *Metrics
	bandwidth     *engine.BandwidthLimiter
	usage         engine.UsageRecorder
	mu            sync.Mutex
	connections   int
	writeMu       sync.Mutex
}

type HandlerOptions struct {
	MaxConnections int
	MaxTunnels     int
	MaxBandwidth   int64
	Heartbeat      time.Duration
	ReadTimeout    time.Duration
	MaxFrameBytes  int64
	DrainTimeout   time.Duration
	Logger         *slog.Logger
	Metrics        *Metrics
	UsageRecorder  engine.UsageRecorder
}

func NewHandler(authenticator AgentAuthenticator, sessions *engine.SessionRegistry, router *engine.RequestRouter, tcp *TCPManager, udp *UDPManager, maxSessions int) (*Handler, error) {
	return NewHandlerWithOptions(authenticator, sessions, router, tcp, udp, HandlerOptions{MaxConnections: maxSessions, MaxTunnels: maxSessions, Heartbeat: 20 * time.Second, ReadTimeout: 90 * time.Second, MaxFrameBytes: 16 << 20})
}

func NewHandlerWithOptions(authenticator AgentAuthenticator, sessions *engine.SessionRegistry, router *engine.RequestRouter, tcp *TCPManager, udp *UDPManager, options HandlerOptions) (*Handler, error) {
	if authenticator == nil || sessions == nil || router == nil || options.MaxConnections < 1 || options.MaxTunnels < 1 || options.Heartbeat <= 0 || options.ReadTimeout <= options.Heartbeat || options.MaxFrameBytes < 1 {
		return nil, fmt.Errorf("authenticator, session registry, router, and positive session limit are required")
	}
	if tcp == nil || udp == nil {
		return nil, fmt.Errorf("tcp and udp managers are required")
	}
	if options.Logger == nil {
		options.Logger = slog.Default()
	}
	if options.Metrics == nil {
		options.Metrics = NewMetrics()
	}
	if options.DrainTimeout < 0 {
		return nil, fmt.Errorf("drain timeout cannot be negative")
	}
	tcp.SetMaxConnections(options.MaxConnections)
	udp.SetMaxPackets(options.MaxConnections)
	return &Handler{authenticator: authenticator, sessions: sessions, router: router, tcp: tcp, udp: udp, maxSessions: options.MaxConnections, maxTunnels: options.MaxTunnels, maxBandwidth: options.MaxBandwidth, heartbeat: options.Heartbeat, readTimeout: options.ReadTimeout, maxFrameBytes: options.MaxFrameBytes, drainTimeout: options.DrainTimeout, logger: options.Logger, metrics: options.Metrics, usage: options.UsageRecorder, bandwidth: engine.NewBandwidthLimiter()}, nil
}

func (h *Handler) Metrics() *Metrics { return h.metrics }

func (h *Handler) Upgrade(c *fiber.Ctx) error {
	return websocket.New(h.Connect)(c)
}

func (h *Handler) Connect(connection *websocket.Conn) {
	if !h.acquireConnection() {
		_ = connection.WriteJSON(protocol.Envelope{Version: protocol.Version, Type: protocol.MessageTypeError, Payload: []byte(`{"code":"capacity","message":"relay capacity reached"}`)})
		_ = connection.Close()
		return
	}
	defer h.releaseConnection()
	h.metrics.AddConnection(1)
	defer h.metrics.AddConnection(-1)
	connection.SetReadLimit(h.maxFrameBytes)
	_ = connection.SetReadDeadline(time.Now().Add(h.readTimeout))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	identity, err := h.authenticator.Authenticate(ctx, bearerToken(connection.Headers("Authorization")))
	if err != nil {
		h.writeJSON(connection, protocol.Envelope{Version: protocol.Version, Type: protocol.MessageTypeError, Payload: []byte(`{"code":"unauthorized","message":"invalid agent token"}`)})
		_ = connection.Close()
		return
	}
	h.logger.Info("relay connection authenticated", slog.String("agent_id", identity.AgentID), slog.String("organization_id", identity.OrganizationID))
	owned := make(map[string]string)
	state := connectionState{}
	connectionCtx, cancelConnection := context.WithCancel(ctx)
	defer cancelConnection()
	go h.sendHeartbeats(connectionCtx, connection, identity.OrganizationID)
	defer func() {
		for tunnelID, sessionID := range owned {
			if h.sessions.Remove(tunnelID, sessionID) {
				h.recordUsage(ctx, identity.OrganizationID, tunnelID, "tunnel_close", 0, 0)
				h.metrics.AddTunnel(-1)
				h.router.RemoveTunnel(tunnelID)
				h.tcp.CloseTunnel(tunnelID)
				h.udp.CloseTunnel(tunnelID)
			}
		}
		h.logger.Info("relay connection closed", slog.String("agent_id", identity.AgentID), slog.Int("tunnels", len(owned)))
	}()
	for {
		messageType, data, err := connection.ReadMessage()
		if err != nil {
			return
		}
		_ = connection.SetReadDeadline(time.Now().Add(h.readTimeout))
		h.metrics.AddFrame(int64(len(data)))
		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			continue
		}
		message, err := protocol.Decode(data)
		if err != nil {
			h.metrics.AddError()
			h.writeError(connection, "protocol", err.Error())
			continue
		}
		if err := h.handleMessage(ctx, connection, identity, message, owned, &state); err != nil {
			h.metrics.AddError()
			h.writeError(connection, "message", err.Error())
		}
		h.recordMessageUsage(ctx, identity.OrganizationID, message)
	}
}

func (h *Handler) sendHeartbeats(ctx context.Context, connection *websocket.Conn, organizationID string) {
	ticker := time.NewTicker(h.heartbeat)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			if err := h.writeJSON(connection, protocol.Envelope{Version: protocol.Version, Type: protocol.MessageTypeHeartbeat, Payload: []byte(fmt.Sprintf(`{"timestamp":%d}`, now.Unix()))}); err != nil {
				_ = connection.Close()
				return
			}
			for _, session := range h.sessions.Snapshot() {
				if session.OrganizationID == organizationID {
					h.sessions.Touch(session.TunnelID)
				}
			}
		}
	}
}

func (h *Handler) handleMessage(ctx context.Context, connection *websocket.Conn, identity AgentIdentity, message protocol.Envelope, owned map[string]string, states ...*connectionState) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	state := &connectionState{negotiated: true, authenticated: true}
	if len(states) > 0 && states[0] != nil {
		state = states[0]
	}
	if h.maxBandwidth > 0 && isDataMessage(message.Type) {
		if err := h.bandwidth.Consume(identity.OrganizationID, h.maxBandwidth, int64(len(message.Payload))); err != nil {
			return err
		}
	}
	switch message.Type {
	case protocol.MessageTypeVersionNegotiate:
		var req protocol.VersionNegotiate
		if err := protocol.DecodePayload(message, &req); err != nil {
			return err
		}
		ack, err := protocol.NegotiateVersion(req)
		if err != nil {
			return err
		}
		state.negotiated = true
		payload, err := protocol.EncodePayload(protocol.MessageTypeVersionNegotiateAck, message.RequestID, ack)
		if err != nil {
			return err
		}
		return h.writeMessage(connection, websocket.TextMessage, payload)
	case protocol.MessageTypeAuth:
		if !state.negotiated {
			return fmt.Errorf("protocol version must be negotiated before authentication")
		}
		var authReq protocol.AuthRequest
		if err := protocol.DecodePayload(message, &authReq); err != nil {
			return err
		}
		id, err := h.authenticator.Authenticate(ctx, authReq.Token)
		if err != nil {
			payload, _ := protocol.EncodePayload(protocol.MessageTypeAuthResponse, message.RequestID, protocol.AuthResponse{Authenticated: false, Error: err.Error()})
			_ = h.writeMessage(connection, websocket.TextMessage, payload)
			return err
		}
		state.authenticated = true
		payload, err := protocol.EncodePayload(protocol.MessageTypeAuthResponse, message.RequestID, protocol.AuthResponse{Authenticated: true, AgentID: id.AgentID, OrganizationID: id.OrganizationID, GrantedCapabilities: []string{"http", "https", "tcp", "udp"}})
		if err != nil {
			return err
		}
		return h.writeMessage(connection, websocket.TextMessage, payload)
	case protocol.MessageTypeFlowControl:
		var fc protocol.FlowControl
		if err := protocol.DecodePayload(message, &fc); err != nil {
			return err
		}
		h.logger.DebugContext(ctx, "flow control message received", slog.String("stream_id", fc.StreamID), slog.String("action", fc.Action))
		return nil
	case protocol.MessageTypeOpenTunnel:
		if !state.negotiated || !state.authenticated {
			return fmt.Errorf("protocol authentication and negotiation are required")
		}
		var open protocol.OpenTunnel
		if err := protocol.DecodePayload(message, &open); err != nil {
			return err
		}
		if open.Protocol == "" || open.LocalPort < 1 || open.LocalPort > 65535 {
			return fmt.Errorf("invalid tunnel open request")
		}
		tunnelID := open.TunnelID
		if tunnelID == "" {
			tunnelID = uuid.NewString()
		}
		previous, exists := h.sessions.Get(tunnelID)
		if exists && previous.OrganizationID != identity.OrganizationID {
			return fmt.Errorf("tunnel belongs to another organization")
		}
		if !exists && len(h.sessions.Snapshot()) >= h.maxTunnels {
			return fmt.Errorf("tunnel capacity reached")
		}
		session := engine.Session{ID: uuid.NewString(), OrganizationID: identity.OrganizationID, TunnelID: tunnelID, Send: func(sendCtx context.Context, outgoing protocol.Envelope) error {
			if err := sendCtx.Err(); err != nil {
				return err
			}
			if h.maxBandwidth > 0 && isDataMessage(outgoing.Type) {
				if err := h.bandwidth.Consume(identity.OrganizationID, h.maxBandwidth, int64(len(outgoing.Payload))); err != nil {
					return err
				}
			}
			h.recordMessageUsage(sendCtx, identity.OrganizationID, outgoing)
			return h.writeJSON(connection, outgoing)
		}, Close: func() { _ = connection.Close() }}
		if err := h.sessions.ReserveWithDrain(session, exists, h.drainTimeout); err != nil {
			return err
		}
		owned[tunnelID] = session.ID
		h.recordUsage(ctx, identity.OrganizationID, tunnelID, "tunnel_open", 0, 1)
		if !exists {
			h.metrics.AddTunnel(1)
		}
		alias := open.CustomDomain
		if alias == "" {
			alias = open.Subdomain
		}
		if alias != "" {
			if err := h.sessions.BindAlias(alias, tunnelID); err != nil {
				h.sessions.Remove(tunnelID, session.ID)
				delete(owned, tunnelID)
				if !exists {
					h.metrics.AddTunnel(-1)
				}
				return err
			}
		}
		publicPort := 0
		var err error
		if open.Protocol == "tcp" {
			if exists {
				h.tcp.SetSender(tunnelID, session.Send)
				publicPort = h.tcp.Port(tunnelID)
			} else {
				publicPort, err = h.tcp.Open(tunnelID, session.Send)
			}
			if err != nil {
				h.sessions.Remove(tunnelID, session.ID)
				delete(owned, tunnelID)
				if !exists {
					h.metrics.AddTunnel(-1)
				}
				return err
			}
		}
		if open.Protocol == "udp" {
			if exists {
				h.udp.SetSender(tunnelID, session.Send)
				publicPort = h.udp.Port(tunnelID)
			} else {
				publicPort, err = h.udp.Open(tunnelID, session.Send)
			}
			if err != nil {
				h.sessions.Remove(tunnelID, session.ID)
				delete(owned, tunnelID)
				if !exists {
					h.metrics.AddTunnel(-1)
				}
				return err
			}
		}
		payload, err := protocol.EncodePayload(protocol.MessageTypeOpenTunnelAck, message.RequestID, protocol.OpenTunnelAck{TunnelID: tunnelID, PublicURL: publicURL(open, tunnelID), PublicPort: publicPort})
		if err != nil {
			h.sessions.Remove(tunnelID, session.ID)
			delete(owned, tunnelID)
			h.tcp.CloseTunnel(tunnelID)
			h.udp.CloseTunnel(tunnelID)
			if !exists {
				h.metrics.AddTunnel(-1)
			}
			return err
		}
		return h.writeMessage(connection, websocket.TextMessage, payload)
	case protocol.MessageTypeHeartbeat:
		for _, session := range h.sessions.Snapshot() {
			if session.OrganizationID == identity.OrganizationID {
				h.sessions.Touch(session.TunnelID)
			}
		}
		return nil
	case protocol.MessageTypeCloseTunnel:
		var closeMessage protocol.CloseTunnel
		if err := protocol.DecodePayload(message, &closeMessage); err != nil {
			return err
		}
		session, ok := h.sessions.Get(closeMessage.TunnelID)
		if !ok || session.OrganizationID != identity.OrganizationID || !h.sessions.Remove(closeMessage.TunnelID, session.ID) {
			return fmt.Errorf("tunnel not found")
		}
		delete(owned, closeMessage.TunnelID)
		h.metrics.AddTunnel(-1)
		h.router.RemoveTunnel(closeMessage.TunnelID)
		h.tcp.CloseTunnel(closeMessage.TunnelID)
		h.udp.CloseTunnel(closeMessage.TunnelID)
		return nil
	case protocol.MessageTypeHTTPResponse:
		if !h.router.Handle(message) {
			return fmt.Errorf("unmatched http response")
		}
		return nil
	case protocol.MessageTypeTCPData:
		var data protocol.TCPData
		if err := protocol.DecodePayload(message, &data); err != nil {
			return err
		}
		if data.TunnelID == "" {
			return fmt.Errorf("tcp tunnel id is required")
		}
		decoded, err := base64.StdEncoding.DecodeString(data.Data)
		if err != nil {
			return fmt.Errorf("decode tcp data: %w", err)
		}
		return h.tcp.Write(data.TunnelID, data.ConnectionID, decoded)
	case protocol.MessageTypeTCPClose:
		var closeMessage protocol.TCPClose
		if err := protocol.DecodePayload(message, &closeMessage); err != nil {
			return err
		}
		if closeMessage.TunnelID == "" {
			return fmt.Errorf("tcp tunnel id is required")
		}
		h.tcp.CloseConnection(closeMessage.TunnelID, closeMessage.ConnectionID)
		return nil
	case protocol.MessageTypeUDPResponse:
		var response protocol.UDPResponse
		if err := protocol.DecodePayload(message, &response); err != nil {
			return err
		}
		return h.udp.Write(response.TunnelID, response)
	default:
		return fmt.Errorf("unsupported relay message type %q", message.Type)
	}
}

func (h *Handler) recordMessageUsage(ctx context.Context, organizationID string, message protocol.Envelope) {
	if h.usage == nil {
		return
	}
	var eventType string
	var encoded string
	switch message.Type {
	case protocol.MessageTypeTCPData:
		var data protocol.TCPData
		if protocol.DecodePayload(message, &data) == nil {
			eventType, encoded = "tcp", data.Data
		}
	case protocol.MessageTypeUDPData:
		var data protocol.UDPData
		if protocol.DecodePayload(message, &data) == nil {
			eventType, encoded = "udp", data.Data
		}
	case protocol.MessageTypeUDPResponse:
		var data protocol.UDPResponse
		if protocol.DecodePayload(message, &data) == nil {
			eventType, encoded = "udp", data.Data
		}
	default:
		return
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return
	}
	tunnelID := tunnelIDFromMessage(message)
	h.recordUsage(ctx, organizationID, tunnelID, eventType, int64(len(data)), 0)
}

func (h *Handler) recordUsage(ctx context.Context, organizationID, tunnelID, eventType string, bytes int64, connections int) {
	if h.usage == nil || organizationID == "" {
		return
	}
	_ = h.usage.Record(ctx, engine.UsageMeasurement{OrganizationID: organizationID, TunnelID: tunnelID, EventType: eventType, Bytes: bytes, Connections: connections})
}

func tunnelIDFromMessage(message protocol.Envelope) string {
	switch message.Type {
	case protocol.MessageTypeTCPData:
		var data protocol.TCPData
		if protocol.DecodePayload(message, &data) == nil {
			return data.TunnelID
		}
	case protocol.MessageTypeUDPData:
		var data protocol.UDPData
		if protocol.DecodePayload(message, &data) == nil {
			return data.TunnelID
		}
	case protocol.MessageTypeUDPResponse:
		var data protocol.UDPResponse
		if protocol.DecodePayload(message, &data) == nil {
			return data.TunnelID
		}
	}
	return ""
}

func (h *Handler) writeError(connection *websocket.Conn, code, message string) {
	payload, err := protocol.EncodePayload(protocol.MessageTypeError, "", protocol.ErrorMessage{Code: code, Message: message})
	if err == nil {
		_ = h.writeMessage(connection, websocket.TextMessage, payload)
	}
}

func (h *Handler) writeJSON(connection *websocket.Conn, value any) error {
	h.writeMu.Lock()
	defer h.writeMu.Unlock()
	_ = connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return connection.WriteJSON(value)
}

func (h *Handler) writeMessage(connection *websocket.Conn, messageType int, data []byte) error {
	h.writeMu.Lock()
	defer h.writeMu.Unlock()
	_ = connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return connection.WriteMessage(messageType, data)
}

func (h *Handler) CloseAll() {
	for _, session := range h.sessions.Snapshot() {
		h.sessions.Remove(session.TunnelID, session.ID)
		h.router.RemoveTunnel(session.TunnelID)
		h.tcp.CloseTunnel(session.TunnelID)
		h.udp.CloseTunnel(session.TunnelID)
	}
}

func isDataMessage(messageType protocol.MessageType) bool {
	return messageType == protocol.MessageTypeHTTPRequest || messageType == protocol.MessageTypeHTTPResponse || messageType == protocol.MessageTypeTCPData || messageType == protocol.MessageTypeUDPData || messageType == protocol.MessageTypeUDPResponse
}

func (h *Handler) acquireConnection() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.connections >= h.maxSessions {
		return false
	}
	h.connections++
	return true
}

func (h *Handler) releaseConnection() {
	h.mu.Lock()
	h.connections--
	h.mu.Unlock()
}
