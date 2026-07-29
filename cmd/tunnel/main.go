package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"codedock.run/codedock-tunnel/internal/config"
	"codedock.run/codedock-tunnel/internal/engine"
	"codedock.run/codedock-tunnel/internal/infra/redis"
	"codedock.run/codedock-tunnel/internal/relay"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var version = "dev"

func main() {
	cfg, err := config.LoadRelay()
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	redisClient, err := redis.Open(ctx, redis.Config{Host: cfg.Redis.Host, Port: cfg.Redis.Port, Password: cfg.Redis.Password, DB: cfg.Redis.DB})
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()
	sessions := engine.NewSessionRegistry()
	requestRouter, err := engine.NewRequestRouter(sessions, cfg.Tunnel.AgentInactivity)
	if err != nil {
		log.Fatal(err)
	}
	httpProxy, err := engine.NewHTTPProxy(cfg.Tunnel.Domain, requestRouter, cfg.Tunnel.MaxBytes)
	if err != nil {
		log.Fatal(err)
	}
	authenticator, err := relay.NewInternalAgentAuthenticator(cfg.Service.InternalAPIURL, cfg.Service.InternalAPISecret, nil)
	if err != nil {
		log.Fatal(err)
	}
	relayHandler, err := relay.NewHandler(authenticator, sessions, cfg.Tunnel.MaxConnections)
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New(fiber.Config{AppName: cfg.App.Name, DisableStartupMessage: true})
	app.Use(recover.New())
	app.Get("/v1/connect", relayHandler.Upgrade)
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "sessions": len(sessions.Snapshot())})
	})
	app.Get("/readyz", func(c *fiber.Ctx) error {
		if err := redisClient.Raw().Ping(c.UserContext()).Err(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "redis unavailable"})
		}
		return c.JSON(fiber.Map{"status": "ready"})
	})
	app.All("/*", adaptor.HTTPHandler(httpProxy))
	go func() {
		if err := app.Listen(cfg.App.ListenAddress()); err != nil {
			log.Printf("relay stopped: %v", err)
			stop()
		}
	}()
	<-ctx.Done()
	if err := app.Shutdown(); err != nil {
		log.Printf("shutdown relay: %v", err)
	}
}
