package relay

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"sync"

	"codedock.run/codedock-tunnel/pkg/protocol"
	"github.com/google/uuid"
)

type UDPManager struct {
	mu        sync.Mutex
	listeners map[string]*net.UDPConn
	packets   map[string]udpPacket
}

type udpPacket struct {
	address  *net.UDPAddr
	listener *net.UDPConn
}

func NewUDPManager() *UDPManager {
	return &UDPManager{listeners: make(map[string]*net.UDPConn), packets: make(map[string]udpPacket)}
}

func (m *UDPManager) Open(tunnelID string, send func(context.Context, protocol.Envelope) error) (int, error) {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero})
	if err != nil {
		return 0, fmt.Errorf("listen public udp port: %w", err)
	}
	m.mu.Lock()
	m.listeners[tunnelID] = listener
	m.mu.Unlock()
	go m.read(tunnelID, listener, send)
	return listener.LocalAddr().(*net.UDPAddr).Port, nil
}

func (m *UDPManager) read(tunnelID string, listener *net.UDPConn, send func(context.Context, protocol.Envelope) error) {
	buffer := make([]byte, 64*1024)
	for {
		count, address, err := listener.ReadFromUDP(buffer)
		if err != nil {
			return
		}
		packetID := uuid.NewString()
		m.mu.Lock()
		m.packets[packetID] = udpPacket{address: address, listener: listener}
		m.mu.Unlock()
		payload, encodeErr := protocol.EncodePayload(protocol.MessageTypeUDPData, "", protocol.UDPData{PacketID: packetID, SourceAddress: address.IP.String(), SourcePort: address.Port, Data: base64.StdEncoding.EncodeToString(buffer[:count])})
		if encodeErr != nil || send(context.Background(), protocol.Envelope{Version: protocol.Version, Type: protocol.MessageTypeUDPData, Payload: payload}) != nil {
			m.mu.Lock()
			delete(m.packets, packetID)
			m.mu.Unlock()
			return
		}
	}
}

func (m *UDPManager) Write(response protocol.UDPResponse) error {
	data, err := base64.StdEncoding.DecodeString(response.Data)
	if err != nil {
		return fmt.Errorf("decode udp data: %w", err)
	}
	m.mu.Lock()
	packet := m.packets[response.PacketID]
	delete(m.packets, response.PacketID)
	m.mu.Unlock()
	if packet.address == nil || packet.listener == nil {
		return fmt.Errorf("udp packet %q not found", response.PacketID)
	}
	_, err = packet.listener.WriteToUDP(data, packet.address)
	return err
}

func (m *UDPManager) CloseTunnel(tunnelID string) {
	m.mu.Lock()
	listener := m.listeners[tunnelID]
	delete(m.listeners, tunnelID)
	if listener != nil {
		for packetID, packet := range m.packets {
			if packet.listener == listener {
				delete(m.packets, packetID)
			}
		}
	}
	m.mu.Unlock()
	if listener != nil {
		_ = listener.Close()
	}
}
