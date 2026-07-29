package protocol

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type FixtureItem struct {
	Name  string `json:"name"`
	Valid bool   `json:"valid"`
	Raw   string `json:"raw"`
}

type FixturesFile struct {
	Fixtures []FixtureItem `json:"fixtures"`
}

func TestProtocolConformance(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "protocol", "fixtures", "conformance_fixtures.json")
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("failed to read conformance fixtures: %v", err)
	}

	var file FixturesFile
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatalf("failed to parse conformance fixtures JSON: %v", err)
	}

	for _, fixture := range file.Fixtures {
		t.Run(fixture.Name, func(t *testing.T) {
			envelope, err := Decode([]byte(fixture.Raw))
			if fixture.Valid {
				if err != nil {
					t.Fatalf("expected valid envelope for %s, got error: %v", fixture.Name, err)
				}
				if envelope.Version < MinSupportedVersion || envelope.Version > MaxSupportedVersion {
					t.Fatalf("expected valid version for %s, got %d", fixture.Name, envelope.Version)
				}
			} else {
				if err == nil {
					t.Fatalf("expected decode error for invalid fixture %s, got nil", fixture.Name)
				}
			}
		})
	}
}

func TestVersionNegotiation(t *testing.T) {
	req := VersionNegotiate{MinVersion: 1, MaxVersion: 1, ClientName: "test", ClientVersion: "1.0"}
	ack, err := NegotiateVersion(req)
	if err != nil {
		t.Fatalf("unexpected negotiation error: %v", err)
	}
	if ack.NegotiatedVersion != 1 {
		t.Fatalf("expected negotiated version 1, got %d", ack.NegotiatedVersion)
	}

	incompatible := VersionNegotiate{MinVersion: 2, MaxVersion: 3}
	if _, err := NegotiateVersion(incompatible); err == nil {
		t.Fatal("expected error negotiating incompatible versions, got nil")
	}
}

func TestMaxFrameSizeEnforcement(t *testing.T) {
	hugeData := make([]byte, AbsoluteMaxFrameSize+100)
	hugeEnvelope := Envelope{Type: MessageTypeData, Version: Version, Payload: hugeData}
	if _, err := Encode(hugeEnvelope); err == nil {
		t.Fatal("expected error encoding frame exceeding max size, got nil")
	}
}
