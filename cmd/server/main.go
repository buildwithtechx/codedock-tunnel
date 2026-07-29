package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"codedock.run/codedock-tunnel/internal/config"
	tunnelhttp "codedock.run/codedock-tunnel/internal/http"
	"codedock.run/codedock-tunnel/internal/infra/postgres"
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
	deps, err := tunnelhttp.NewDatabaseDependencies(db, cfg.Auth)
	if err != nil {
		log.Fatal(err)
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
