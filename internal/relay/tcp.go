package relay

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"sync"

	"codedock.run/codedock-tunnel/pkg/protocol"
	"github.com/google/uuid"
)

type TCPManager struct {
	mu          sync.Mutex
	listeners   map[string]net.Listener
	connections map[string]net.Conn
	tunnels     map[string]map[string]struct{}
	senders     map[string]func(context.Context, protocol.Envelope) error
	max         int
}

func NewTCPManager() *TCPManager {
	return &TCPManager{listeners: make(map[string]net.Listener), connections: make(map[string]net.Conn), tunnels: make(map[string]map[string]struct{}), senders: make(map[string]func(context.Context, protocol.Envelope) error), max: 1000}
}

func (m *TCPManager) SetMaxConnections(max int) {
	if max < 1 {
		return
	}
	m.mu.Lock()
	m.max = max
	m.mu.Unlock()
}

func (m *TCPManager) Open(tunnelID string, send func(context.Context, protocol.Envelope) error) (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("listen public tcp port: %w", err)
	}
	m.mu.Lock()
	m.listeners[tunnelID] = listener
	m.tunnels[tunnelID] = make(map[string]struct{})
	m.senders[tunnelID] = send
	m.mu.Unlock()
	go m.accept(tunnelID, listener)
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func (m *TCPManager) accept(tunnelID string, listener net.Listener) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			return
		}
		m.mu.Lock()
		atCapacity := len(m.connections) >= m.max
		m.mu.Unlock()
		if atCapacity {
			_ = connection.Close()
			continue
		}
		connectionID := uuid.NewString()
		m.mu.Lock()
		m.connections[connectionID] = connection
		m.tunnels[tunnelID][connectionID] = struct{}{}
		m.mu.Unlock()
		go m.read(tunnelID, connectionID, connection)
	}
}

func (m *TCPManager) read(tunnelID, connectionID string, connection net.Conn) {
	defer m.CloseConnection(tunnelID, connectionID)
	buffer := make([]byte, 32*1024)
	for {
		count, err := connection.Read(buffer)
		if count > 0 {
			send := m.sender(tunnelID)
			if send == nil {
				return
			}
			payload, encodeErr := protocol.EncodePayload(protocol.MessageTypeTCPData, "", protocol.TCPData{TunnelID: tunnelID, ConnectionID: connectionID, Data: base64.StdEncoding.EncodeToString(buffer[:count])})
			if encodeErr != nil || send(context.Background(), protocol.Envelope{Version: protocol.Version, Type: protocol.MessageTypeTCPData, Payload: payload}) != nil {
				return
			}
		}
		if err != nil {
			if err != io.EOF {
				send := m.sender(tunnelID)
				if send == nil {
					return
				}
				payload, _ := protocol.EncodePayload(protocol.MessageTypeTCPClose, "", protocol.TCPClose{TunnelID: tunnelID, ConnectionID: connectionID})
				_ = send(context.Background(), protocol.Envelope{Version: protocol.Version, Type: protocol.MessageTypeTCPClose, Payload: payload})
			}
			return
		}
	}
}

func (m *TCPManager) SetSender(tunnelID string, send func(context.Context, protocol.Envelope) error) {
	m.mu.Lock()
	if _, ok := m.listeners[tunnelID]; ok {
		m.senders[tunnelID] = send
	}
	m.mu.Unlock()
}

func (m *TCPManager) Port(tunnelID string) int {
	m.mu.Lock()
	listener := m.listeners[tunnelID]
	m.mu.Unlock()
	if listener == nil {
		return 0
	}
	return listener.Addr().(*net.TCPAddr).Port
}

func (m *TCPManager) sender(tunnelID string) func(context.Context, protocol.Envelope) error {
	m.mu.Lock()
	send := m.senders[tunnelID]
	m.mu.Unlock()
	return send
}

func (m *TCPManager) Write(tunnelID, connectionID string, data []byte) error {
	m.mu.Lock()
	connection, ok := m.connections[connectionID]
	owned := false
	if tunnelConnections := m.tunnels[tunnelID]; tunnelConnections != nil {
		_, owned = tunnelConnections[connectionID]
	}
	m.mu.Unlock()
	if !ok || !owned {
		return fmt.Errorf("tcp connection %q not found", connectionID)
	}
	_, err := connection.Write(data)
	return err
}

func (m *TCPManager) CloseConnection(tunnelID, connectionID string) {
	m.mu.Lock()
	connection := m.connections[connectionID]
	if tunnelConnections := m.tunnels[tunnelID]; tunnelConnections != nil {
		if _, ok := tunnelConnections[connectionID]; !ok {
			m.mu.Unlock()
			return
		}
	}
	delete(m.connections, connectionID)
	for tunnelID, connections := range m.tunnels {
		delete(connections, connectionID)
		if len(connections) == 0 && m.listeners[tunnelID] == nil {
			delete(m.tunnels, tunnelID)
		}
	}
	m.mu.Unlock()
	if connection != nil {
		_ = connection.Close()
	}
}

func (m *TCPManager) CloseTunnel(tunnelID string) {
	m.mu.Lock()
	listener := m.listeners[tunnelID]
	delete(m.listeners, tunnelID)
	delete(m.senders, tunnelID)
	connections := m.tunnels[tunnelID]
	delete(m.tunnels, tunnelID)
	for connectionID := range connections {
		if connection := m.connections[connectionID]; connection != nil {
			_ = connection.Close()
		}
		delete(m.connections, connectionID)
	}
	m.mu.Unlock()
	if listener != nil {
		_ = listener.Close()
	}
}
