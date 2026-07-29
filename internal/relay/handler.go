package relay

import (
	"context"
	"fmt"
	"strings"
	"sync"

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

type Handler struct {
	authenticator AgentAuthenticator
	sessions      *engine.SessionRegistry
	maxSessions   int
	mu            sync.Mutex
	connections   int
}

func NewHandler(authenticator AgentAuthenticator, sessions *engine.SessionRegistry, maxSessions int) (*Handler, error) {
	if authenticator == nil || sessions == nil || maxSessions < 1 {
		return nil, fmt.Errorf("authenticator, session registry, and positive session limit are required")
	}
	return &Handler{authenticator: authenticator, sessions: sessions, maxSessions: maxSessions}, nil
}

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	identity, err := h.authenticator.Authenticate(ctx, bearerToken(connection.Headers("Authorization")))
	if err != nil {
		_ = connection.WriteJSON(protocol.Envelope{Version: protocol.Version, Type: protocol.MessageTypeError, Payload: []byte(`{"code":"unauthorized","message":"invalid agent token"}`)})
		_ = connection.Close()
		return
	}
	for {
		messageType, data, err := connection.ReadMessage()
		if err != nil {
			return
		}
		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			continue
		}
		message, err := protocol.Decode(data)
		if err != nil {
			h.writeError(connection, "protocol", err.Error())
			continue
		}
		if err := h.handleMessage(ctx, connection, identity, message); err != nil {
			h.writeError(connection, "message", err.Error())
		}
	}
}

func (h *Handler) handleMessage(ctx context.Context, connection *websocket.Conn, identity AgentIdentity, message protocol.Envelope) error {
	switch message.Type {
	case protocol.MessageTypeOpenTunnel:
		var open protocol.OpenTunnel
		if err := protocol.DecodePayload(message, &open); err != nil {
			return err
		}
		if open.Protocol == "" || open.LocalPort < 1 || open.LocalPort > 65535 {
			return fmt.Errorf("invalid tunnel open request")
		}
		tunnelID := uuid.NewString()
		session := engine.Session{ID: uuid.NewString(), OrganizationID: identity.OrganizationID, TunnelID: tunnelID, Send: func(_ context.Context, outgoing protocol.Envelope) error {
			return connection.WriteJSON(outgoing)
		}, Close: func() { _ = connection.Close() }}
		if err := h.sessions.Reserve(session, false); err != nil {
			return err
		}
		payload, err := protocol.EncodePayload(protocol.MessageTypeOpenTunnelAck, message.RequestID, protocol.OpenTunnelAck{TunnelID: tunnelID, PublicURL: publicURL(open, tunnelID)})
		if err != nil {
			return err
		}
		return connection.WriteMessage(websocket.TextMessage, payload)
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
		if !h.sessions.Remove(closeMessage.TunnelID, "") {
			return fmt.Errorf("tunnel not found")
		}
		return nil
	default:
		return fmt.Errorf("unsupported relay message type %q", message.Type)
	}
}

func (h *Handler) writeError(connection *websocket.Conn, code, message string) {
	payload, err := protocol.EncodePayload(protocol.MessageTypeError, "", protocol.ErrorMessage{Code: code, Message: message})
	if err == nil {
		_ = connection.WriteMessage(websocket.TextMessage, payload)
	}
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

func bearerToken(value string) string {
	value = strings.TrimSpace(value)
	parts := strings.SplitN(value, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func publicURL(open protocol.OpenTunnel, tunnelID string) string {
	if open.CustomDomain != "" {
		return "https://" + strings.TrimSuffix(open.CustomDomain, ".")
	}
	name := open.Subdomain
	if name == "" {
		name = strings.Split(tunnelID, "-")[0]
	}
	return "https://" + name + ".tunnel.codedock-tunnel.dev"
}
