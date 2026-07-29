package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codedock.run/codedock-tunnel/internal/config"
	"codedock.run/codedock-tunnel/pkg/client"
	"codedock.run/codedock-tunnel/pkg/protocol"
)

var version = "dev"

func main() {
	cfg, err := config.LoadCLI()
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) < 2 || os.Args[1] == "help" {
		printUsage()
		return
	}
	if os.Args[1] == "version" {
		fmt.Println(version)
		return
	}
	if os.Args[1] != "health" {
		if os.Args[1] != "open" {
			log.Fatalf("unknown command %q", os.Args[1])
		}
		openTunnel(cfg)
		return
	}
	flags := flag.NewFlagSet("health", flag.ExitOnError)
	apiURL := flags.String("api-url", cfg.APIURL, "tunnel API URL")
	_ = flags.Parse(os.Args[2:])
	apiClient, err := client.New(client.Config{BaseURL: *apiURL, APIKey: cfg.APIKey})
	if err != nil {
		log.Fatal(err)
	}
	if err := apiClient.Do(context.Background(), http.MethodGet, "/readyz", nil, nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("ready")
}

func openTunnel(cfg config.CLIConfig) {
	flags := flag.NewFlagSet("open", flag.ExitOnError)
	port := flags.Int("port", 3000, "local port")
	protocolName := flags.String("protocol", "http", "tunnel protocol")
	subdomain := flags.String("subdomain", "", "requested subdomain")
	_ = flags.Parse(os.Args[2:])
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	connection, err := client.OpenRelay(ctx, client.RelayConfig{URL: cfg.RelayURL, Token: cfg.AgentToken}, protocol.OpenTunnel{LocalPort: *port, Protocol: *protocolName, Subdomain: *subdomain})
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	fmt.Printf("tunnel %s %s\n", connection.TunnelID, connection.PublicURL)
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			if err := connection.SendHeartbeat(); err != nil {
				return
			}
		}
	}()
	if err := connection.ServeLocal(ctx, "http://"+net.JoinHostPort("127.0.0.1", fmt.Sprint(*port))); err != nil && ctx.Err() == nil {
		log.Fatal(err)
	}
}

func printUsage() {
	fmt.Println("codedock-tunnel <command>")
	fmt.Println("commands: health, open, version")
}
