package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
	target := "http://127.0.0.1:" + fmt.Sprint(*port)
	if *protocolName == "tcp" || *protocolName == "udp" {
		target = "127.0.0.1:" + fmt.Sprint(*port)
	}
	delay := 2 * time.Second
	for ctx.Err() == nil {
		connection, err := client.OpenRelay(ctx, client.RelayConfig{URL: cfg.RelayURL, Token: cfg.AgentToken}, protocol.OpenTunnel{LocalPort: *port, Protocol: *protocolName, Subdomain: *subdomain})
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("relay connection failed: %v; retrying in %s", err, delay)
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
			}
			if delay < 30*time.Second {
				delay *= 2
			}
			continue
		}
		delay = 2 * time.Second
		if connection.PublicPort > 0 {
			fmt.Printf("tunnel %s %s:%d\n", connection.TunnelID, connection.PublicURL, connection.PublicPort)
		} else {
			fmt.Printf("tunnel %s %s\n", connection.TunnelID, connection.PublicURL)
		}
		ticker := time.NewTicker(20 * time.Second)
		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err := connection.SendHeartbeat(); err != nil {
						return
					}
				}
			}
		}()
		serveErr := connection.ServeLocal(ctx, target)
		ticker.Stop()
		connection.Close()
		<-done
		if serveErr != nil && ctx.Err() == nil {
			log.Printf("relay connection closed: %v; reconnecting", serveErr)
		}
	}
}

func printUsage() {
	fmt.Println("codedock-tunnel <command>")
	fmt.Println("commands: health, open, version")
}
