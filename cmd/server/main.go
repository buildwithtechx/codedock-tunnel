package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"codedock.run/codedock-tunnel/internal/auth"
	"codedock.run/codedock-tunnel/internal/config"
	tunnelhttp "codedock.run/codedock-tunnel/internal/http"
	"codedock.run/codedock-tunnel/internal/infra/postgres"
	"codedock.run/codedock-tunnel/internal/infra/redis"
	"codedock.run/codedock-tunnel/internal/infra/storage"
	"codedock.run/codedock-tunnel/internal/services"
)

func main() {
	cfg, err := config.LoadAPI()
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	db, err := postgres.Open(ctx, postgres.Config{DSN: cfg.Database.URL, MaxOpenConns: cfg.Database.MaxConns, MaxIdleConns: cfg.Database.MaxConns, ConnMaxLifetime: cfg.Database.MaxLifetime})
	if err != nil {
		log.Fatal(err)
	}
	if err := postgres.Migrate(db); err != nil {
		log.Fatal(err)
	}
	deps, err := tunnelhttp.NewDatabaseDependencies(db, cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Auth.EncryptionKey != "" {
		secretBox, err := storage.NewSecretBox([]byte(cfg.Auth.EncryptionKey))
		if err != nil {
			log.Fatal(err)
		}
		deps.Auth.SetSecretProtector(secretBox)
		deps.DeviceLogin.SetSecretProtector(secretBox)
		deps.Billing.SetSecretProtector(secretBox)
		deps.Billing.SetBillingSecretProtector(secretBox)
		deps.Billing.SetGracePeriod(cfg.Billing.GracePeriod)
	}
	redisClient, err := redis.Open(ctx, redis.Config{Host: cfg.Redis.Host, Port: cfg.Redis.Port, Password: cfg.Redis.Password, DB: cfg.Redis.DB})
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()
	providers := auth.NewOAuthProviders(auth.OAuthConfig{GoogleClientID: cfg.Auth.GoogleClientID, GoogleClientSecret: cfg.Auth.GoogleClientSecret, GitHubClientID: cfg.Auth.GitHubClientID, GitHubClientSecret: cfg.Auth.GitHubClientSecret})
	if len(providers) > 0 {
		stateStore, err := redis.NewOAuthStateStore(redisClient, cfg.Auth.OAuthStateTTL)
		if err != nil {
			log.Fatal(err)
		}
		oauthService, err := services.NewOAuthService(deps.Auth, providers, stateStore)
		if err != nil {
			log.Fatal(err)
		}
		deps.OAuth = oauthService
	}
	server, err := tunnelhttp.NewServer(cfg, deps)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := server.Listen(cfg.App.ListenAddress()); err != nil {
			log.Printf("http server stopped: %v", err)
			stop()
		}
	}()
	<-ctx.Done()
	if err := server.Shutdown(); err != nil {
		log.Printf("shutdown http server: %v", err)
	}
}
