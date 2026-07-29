package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
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
	conn       *websocket.Conn
	TunnelID   string
	PublicURL  string
	PublicPort int
	Protocol   string
	writeMu    sync.Mutex
	tcpMu      sync.Mutex
	tcpConns   map[string]net.Conn
	udpMu      sync.Mutex
	udpConn    *net.UDPConn
	udpQueue   []string
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
	return &RelayConnection{conn: connection, TunnelID: ack.TunnelID, PublicURL: ack.PublicURL, PublicPort: ack.PublicPort, Protocol: open.Protocol, tcpConns: make(map[string]net.Conn)}, nil
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
		switch message.Type {
		case protocol.MessageTypeHTTPRequest:
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
		case protocol.MessageTypeTCPData:
			if err := c.handleTCPData(targetURL, message); err != nil {
				return err
			}
		case protocol.MessageTypeTCPClose:
			var closeMessage protocol.TCPClose
			if err := protocol.DecodePayload(message, &closeMessage); err != nil {
				return err
			}
			c.closeTCP(closeMessage.ConnectionID)
		case protocol.MessageTypeUDPData:
			if err := c.handleUDPData(targetURL, message); err != nil {
				return err
			}
		}
	}
}

func (c *RelayConnection) handleUDPData(target string, message protocol.Envelope) error {
	var data protocol.UDPData
	if err := protocol.DecodePayload(message, &data); err != nil {
		return err
	}
	decoded, err := base64.StdEncoding.DecodeString(data.Data)
	if err != nil {
		return err
	}
	c.udpMu.Lock()
	created := c.udpConn == nil
	if created {
		address, resolveErr := net.ResolveUDPAddr("udp", target)
		if resolveErr != nil {
			c.udpMu.Unlock()
			return resolveErr
		}
		c.udpConn, err = net.DialUDP("udp", nil, address)
	}
	if err == nil {
		c.udpQueue = append(c.udpQueue, data.PacketID)
	}
	connection := c.udpConn
	c.udpMu.Unlock()
	if err != nil {
		return err
	}
	if created {
		go c.proxyUDP(connection)
	}
	_, err = connection.Write(decoded)
	return err
}

func (c *RelayConnection) proxyUDP(connection *net.UDPConn) {
	buffer := make([]byte, 64*1024)
	for {
		count, address, err := connection.ReadFromUDP(buffer)
		if err != nil {
			return
		}
		c.udpMu.Lock()
		if len(c.udpQueue) == 0 {
			c.udpMu.Unlock()
			continue
		}
		packetID := c.udpQueue[0]
		c.udpQueue = c.udpQueue[1:]
		c.udpMu.Unlock()
		payload, encodeErr := protocol.EncodePayload(protocol.MessageTypeUDPResponse, "", protocol.UDPResponse{PacketID: packetID, TargetAddress: address.IP.String(), TargetPort: address.Port, Data: base64.StdEncoding.EncodeToString(buffer[:count])})
		if encodeErr != nil || c.write(payload) != nil {
			return
		}
	}
}

func (c *RelayConnection) handleTCPData(target string, message protocol.Envelope) error {
	var data protocol.TCPData
	if err := protocol.DecodePayload(message, &data); err != nil {
		return err
	}
	c.tcpMu.Lock()
	connection := c.tcpConns[data.ConnectionID]
	created := false
	var err error
	if connection == nil {
		connection, err = net.Dial("tcp", target)
		if err == nil {
			c.tcpConns[data.ConnectionID] = connection
			created = true
		}
	}
	c.tcpMu.Unlock()
	if connection == nil {
		return c.sendTCPClose(data.ConnectionID, "connect local target failed")
	}

	decoded, err := base64.StdEncoding.DecodeString(data.Data)
	if err != nil {
		c.closeTCP(data.ConnectionID)
		return err
	}
	if _, err := connection.Write(decoded); err != nil {
		c.closeTCP(data.ConnectionID)
		return err
	}
	if created {
		go c.proxyTCP(data.ConnectionID, connection)
	}
	return nil
}

func (c *RelayConnection) proxyTCP(connectionID string, connection net.Conn) {
	defer connection.Close()
	buffer := make([]byte, 32*1024)
	for {
		count, err := connection.Read(buffer)
		if count > 0 {
			payload, encodeErr := protocol.EncodePayload(protocol.MessageTypeTCPData, "", protocol.TCPData{ConnectionID: connectionID, Data: base64.StdEncoding.EncodeToString(buffer[:count])})
			if encodeErr != nil || c.write(payload) != nil {
				return
			}
		}
		if err != nil {
			c.closeTCP(connectionID)
			_ = c.sendTCPClose(connectionID, err.Error())
			return
		}
	}
}

func (c *RelayConnection) closeTCP(connectionID string) {
	c.tcpMu.Lock()
	connection := c.tcpConns[connectionID]
	delete(c.tcpConns, connectionID)
	c.tcpMu.Unlock()
	if connection != nil {
		_ = connection.Close()
	}
}

func (c *RelayConnection) sendTCPClose(connectionID, reason string) error {
	payload, err := protocol.EncodePayload(protocol.MessageTypeTCPClose, "", protocol.TCPClose{ConnectionID: connectionID, Reason: reason})
	if err != nil {
		return err
	}
	return c.write(payload)
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
	c.tcpMu.Lock()
	for connectionID, connection := range c.tcpConns {
		_ = connection.Close()
		delete(c.tcpConns, connectionID)
	}
	c.tcpMu.Unlock()
	c.udpMu.Lock()
	if c.udpConn != nil {
		_ = c.udpConn.Close()
		c.udpConn = nil
	}
	c.udpQueue = nil
	c.udpMu.Unlock()
	return c.conn.Close()
}
