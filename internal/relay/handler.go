package relay

import (
	"context"
	"encoding/base64"
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
	router        *engine.RequestRouter
	tcp           *TCPManager
	udp           *UDPManager
	maxSessions   int
	mu            sync.Mutex
	connections   int
	writeMu       sync.Mutex
}

func NewHandler(authenticator AgentAuthenticator, sessions *engine.SessionRegistry, router *engine.RequestRouter, tcp *TCPManager, udp *UDPManager, maxSessions int) (*Handler, error) {
	if authenticator == nil || sessions == nil || router == nil || maxSessions < 1 {
		return nil, fmt.Errorf("authenticator, session registry, router, and positive session limit are required")
	}
	if tcp == nil || udp == nil {
		return nil, fmt.Errorf("tcp and udp managers are required")
	}
	return &Handler{authenticator: authenticator, sessions: sessions, router: router, tcp: tcp, udp: udp, maxSessions: maxSessions}, nil
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
		h.writeJSON(connection, protocol.Envelope{Version: protocol.Version, Type: protocol.MessageTypeError, Payload: []byte(`{"code":"unauthorized","message":"invalid agent token"}`)})
		_ = connection.Close()
		return
	}
	owned := make(map[string]string)
	defer func() {
		for tunnelID, sessionID := range owned {
			h.sessions.Remove(tunnelID, sessionID)
			h.router.RemoveTunnel(tunnelID)
			h.tcp.CloseTunnel(tunnelID)
			h.udp.CloseTunnel(tunnelID)
		}
	}()
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
		if err := h.handleMessage(ctx, connection, identity, message, owned); err != nil {
			h.writeError(connection, "message", err.Error())
		}
	}
}

func (h *Handler) handleMessage(ctx context.Context, connection *websocket.Conn, identity AgentIdentity, message protocol.Envelope, owned map[string]string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
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
		session := engine.Session{ID: uuid.NewString(), OrganizationID: identity.OrganizationID, TunnelID: tunnelID, Send: func(sendCtx context.Context, outgoing protocol.Envelope) error {
			if err := sendCtx.Err(); err != nil {
				return err
			}
			return h.writeJSON(connection, outgoing)
		}, Close: func() { _ = connection.Close() }}
		if err := h.sessions.Reserve(session, false); err != nil {
			return err
		}
		owned[tunnelID] = session.ID
		publicPort := 0
		var err error
		if open.Protocol == "tcp" {
			publicPort, err = h.tcp.Open(tunnelID, session.Send)
			if err != nil {
				h.sessions.Remove(tunnelID, session.ID)
				return err
			}
		}
		if open.Protocol == "udp" {
			publicPort, err = h.udp.Open(tunnelID, session.Send)
			if err != nil {
				h.sessions.Remove(tunnelID, session.ID)
				return err
			}
		}
		payload, err := protocol.EncodePayload(protocol.MessageTypeOpenTunnelAck, message.RequestID, protocol.OpenTunnelAck{TunnelID: tunnelID, PublicURL: publicURL(open, tunnelID), PublicPort: publicPort})
		if err != nil {
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
		if !h.sessions.Remove(closeMessage.TunnelID, "") {
			return fmt.Errorf("tunnel not found")
		}
		delete(owned, closeMessage.TunnelID)
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
		decoded, err := base64.StdEncoding.DecodeString(data.Data)
		if err != nil {
			return fmt.Errorf("decode tcp data: %w", err)
		}
		return h.tcp.Write(data.ConnectionID, decoded)
	case protocol.MessageTypeTCPClose:
		var closeMessage protocol.TCPClose
		if err := protocol.DecodePayload(message, &closeMessage); err != nil {
			return err
		}
		h.tcp.CloseConnection(closeMessage.ConnectionID)
		return nil
	case protocol.MessageTypeUDPResponse:
		var response protocol.UDPResponse
		if err := protocol.DecodePayload(message, &response); err != nil {
			return err
		}
		return h.udp.Write(response)
	default:
		return fmt.Errorf("unsupported relay message type %q", message.Type)
	}
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
	return connection.WriteJSON(value)
}

func (h *Handler) writeMessage(connection *websocket.Conn, messageType int, data []byte) error {
	h.writeMu.Lock()
	defer h.writeMu.Unlock()
	return connection.WriteMessage(messageType, data)
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
