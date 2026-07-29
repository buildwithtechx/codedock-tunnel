package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"codedock.run/codedock-tunnel/internal/auth"
	"codedock.run/codedock-tunnel/internal/models"
	"codedock.run/codedock-tunnel/internal/repositories"
)

type DeviceLoginService struct {
	logins repositories.DeviceLoginRepository
	now    func() time.Time
	ttl    time.Duration
}

func NewDeviceLoginService(logins repositories.DeviceLoginRepository, ttl time.Duration) (*DeviceLoginService, error) {
	if logins == nil {
		return nil, fmt.Errorf("device login repository is required")
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("device login ttl must be positive")
	}
	return &DeviceLoginService{logins: logins, now: time.Now, ttl: ttl}, nil
}

func (s *DeviceLoginService) Start(ctx context.Context, ipAddress string) (string, models.DeviceLogin, error) {
	code, err := auth.NewToken("", 24)
	if err != nil {
		return "", models.DeviceLogin{}, err
	}
	now := s.now()
	login := models.DeviceLogin{CodeHash: auth.HashToken(code), Status: "pending", ExpiresAt: now.Add(s.ttl), IPAddress: ipAddress}
	if err := s.logins.Create(ctx, &login); err != nil {
		return "", models.DeviceLogin{}, fmt.Errorf("create device login: %w", err)
	}
	return code, login, nil
}

func (s *DeviceLoginService) Complete(ctx context.Context, code, userID string) (string, error) {
	if strings.TrimSpace(code) == "" || userID == "" {
		return "", fmt.Errorf("device code and user id are required")
	}
	login, err := s.logins.FindPending(ctx, auth.HashToken(code), s.now())
	if err != nil {
		return "", fmt.Errorf("find device login: %w", err)
	}
	token, err := auth.NewToken("cdt", 32)
	if err != nil {
		return "", err
	}
	if err := s.logins.Complete(ctx, login.ID, userID, auth.HashToken(token), s.now()); err != nil {
		return "", fmt.Errorf("complete device login: %w", err)
	}
	return token, nil
}

func (s *DeviceLoginService) Consume(ctx context.Context, token string) (models.DeviceLogin, error) {
	if strings.TrimSpace(token) == "" {
		return models.DeviceLogin{}, fmt.Errorf("device token is required")
	}
	login, err := s.logins.ConsumeToken(ctx, auth.HashToken(token), s.now())
	if err != nil {
		return models.DeviceLogin{}, fmt.Errorf("consume device login: %w", err)
	}
	return login, nil
}
