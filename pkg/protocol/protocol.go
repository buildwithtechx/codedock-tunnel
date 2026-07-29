package protocol

import (
	"encoding/json"
	"fmt"

	generated "codedock.run/codedock-tunnel/protocol/generated/go"
)

const Version = generated.Version

type MessageType = generated.MessageType
type Envelope = generated.Envelope
type OpenTunnel = generated.OpenTunnel
type OpenTunnelAck = generated.OpenTunnelAck
type CloseTunnel = generated.CloseTunnel
type Data = generated.Data
type Heartbeat = generated.Heartbeat
type ErrorMessage = generated.ErrorMessage

const (
	MessageTypeOpenTunnel    = generated.MessageTypeOpenTunnel
	MessageTypeOpenTunnelAck = generated.MessageTypeOpenTunnelAck
	MessageTypeCloseTunnel   = generated.MessageTypeCloseTunnel
	MessageTypeData          = generated.MessageTypeData
	MessageTypeHeartbeat     = generated.MessageTypeHeartbeat
	MessageTypeError         = generated.MessageTypeError
)

func Encode(message Envelope) ([]byte, error) {
	if message.Version == 0 {
		message.Version = Version
	}
	if message.Type == "" {
		return nil, fmt.Errorf("protocol message type is required")
	}
	return json.Marshal(message)
}

func Decode(data []byte) (Envelope, error) {
	var message Envelope
	if err := json.Unmarshal(data, &message); err != nil {
		return Envelope{}, fmt.Errorf("decode protocol message: %w", err)
	}
	if message.Version != Version {
		return Envelope{}, fmt.Errorf("unsupported protocol version %d", message.Version)
	}
	if message.Type == "" {
		return Envelope{}, fmt.Errorf("protocol message type is required")
	}
	return message, nil
}

func EncodePayload(messageType MessageType, requestID string, payload any) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode protocol payload: %w", err)
	}
	return Encode(Envelope{Type: messageType, RequestID: requestID, Payload: body})
}

func DecodePayload(message Envelope, target any) error {
	if err := json.Unmarshal(message.Payload, target); err != nil {
		return fmt.Errorf("decode protocol payload: %w", err)
	}
	return nil
}
