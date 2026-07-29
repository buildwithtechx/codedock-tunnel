package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"codedock.run/codedock-tunnel/internal/config"
	"codedock.run/codedock-tunnel/pkg/client"
)

type DeviceCodeResponse struct {
	DeviceCode      string `json:"deviceCode"`
	UserCode        string `json:"userCode"`
	VerificationURL string `json:"verificationUrl"`
	ExpiresIn       int    `json:"expiresIn"`
	Interval        int    `json:"interval"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
}

func runLogin(cfg config.CLIConfig) {
	apiClient, err := client.New(client.Config{BaseURL: cfg.APIURL})
	if err != nil {
		log.Fatalf("initialize client: %v", err)
	}

	var codeResp DeviceCodeResponse
	if err := apiClient.Do(context.Background(), http.MethodPost, "/api/v1/auth/device/code", nil, &codeResp); err != nil {
		log.Fatalf("initiate device login: %v", err)
	}

	fmt.Printf("User Code: %s\nVerification URL: %s\nWaiting for authorization...\n", codeResp.UserCode, codeResp.VerificationURL)

	interval := time.Duration(codeResp.Interval) * time.Second
	if interval < time.Second {
		interval = 5 * time.Second
	}
	deadline := time.Now().Add(time.Duration(codeResp.ExpiresIn) * time.Second)

	for time.Now().Before(deadline) {
		time.Sleep(interval)
		var tokenResp TokenResponse
		err := apiClient.Do(context.Background(), http.MethodPost, "/api/v1/auth/device/token", map[string]string{"deviceCode": codeResp.DeviceCode}, &tokenResp)
		if err == nil && tokenResp.AccessToken != "" {
			cfg.APIKey = tokenResp.AccessToken
			if err := config.SaveCLI(cfg); err != nil {
				log.Fatalf("save credentials: %v", err)
			}
			fmt.Println("Successfully logged in!")
			return
		}
	}
	log.Fatal("device login expired")
}
