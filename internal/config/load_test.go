package config

import "testing"

func TestLoadUsesCodedockEnvironmentPrefix(t *testing.T) {
	t.Setenv("CODEDOCK_PORT", "9090")
	t.Setenv("CODEDOCK_APP_NAME", "test-tunnel")
	t.Setenv("CODEDOCK_DATABASE_MAX_CONNS", "12")
	t.Setenv("CODEDOCK_GOOGLE_CLIENT_ID", "google-client")
	t.Setenv("CODEDOCK_GITHUB_CLIENT_ID", "github-client")
	t.Setenv("CODEDOCK_ZEPTO_API_KEY", "zepto-key")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.App.Port != "9090" {
		t.Fatalf("expected port 9090, got %s", cfg.App.Port)
	}
	if cfg.App.Name != "test-tunnel" {
		t.Fatalf("expected test-tunnel, got %s", cfg.App.Name)
	}
	if cfg.Database.MaxConns != 12 {
		t.Fatalf("expected 12 database connections, got %d", cfg.Database.MaxConns)
	}
	if cfg.Auth.GoogleClientID != "google-client" || cfg.Auth.GitHubClientID != "github-client" {
		t.Fatal("oauth configuration did not load")
	}
	if cfg.Mail.ZeptoAPIKey != "zepto-key" {
		t.Fatal("zepto configuration did not load")
	}
}
