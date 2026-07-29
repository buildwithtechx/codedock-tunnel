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
	"github.com/gofiber/fiber/v2"
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
	if _, err := engine.NewRequestRouter(sessions, cfg.Tunnel.AgentInactivity); err != nil {
		log.Fatal(err)
	}
	app := fiber.New(fiber.Config{AppName: cfg.App.Name, DisableStartupMessage: true})
	app.Use(recover.New())
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "sessions": len(sessions.Snapshot())})
	})
	app.Get("/readyz", func(c *fiber.Ctx) error {
		if err := redisClient.Raw().Ping(c.UserContext()).Err(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "redis unavailable"})
		}
		return c.JSON(fiber.Map{"status": "ready"})
	})
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
