package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"codedock.run/codedock-tunnel/internal/config"
	"codedock.run/codedock-tunnel/pkg/client"
)

type TunnelDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Protocol  string `json:"protocol"`
	Subdomain string `json:"subdomain"`
	Status    string `json:"status"`
	PublicURL string `json:"publicUrl"`
}

func printOutput(jsonOutput bool, value any) {
	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(value)
		return
	}
	switch v := value.(type) {
	case []TunnelDTO:
		for _, t := range v {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", t.ID, t.Name, t.Protocol, t.Status, t.PublicURL)
		}
	case TunnelDTO:
		fmt.Printf("ID: %s\nName: %s\nProtocol: %s\nStatus: %s\nPublic URL: %s\n", v.ID, v.Name, v.Protocol, v.Status, v.PublicURL)
	default:
		fmt.Println(v)
	}
}

func runTunnelsCommand(cfg config.CLIConfig, args []string) {
	if len(args) == 0 {
		log.Fatal("tunnel action required: create, list, inspect, start, stop, revoke")
	}
	action := args[0]
	flags := flag.NewFlagSet(action, flag.ExitOnError)
	jsonOutput := flags.Bool("json", false, "output in JSON format")
	_ = flags.Parse(args[1:])

	apiClient, err := client.New(client.Config{BaseURL: cfg.APIURL, APIKey: cfg.APIKey})
	if err != nil {
		log.Fatalf("initialize client: %v", err)
	}

	switch action {
	case "list":
		var tunnels []TunnelDTO
		if err := apiClient.Do(context.Background(), http.MethodGet, "/api/v1/tunnels", nil, &tunnels); err != nil {
			log.Fatalf("list tunnels: %v", err)
		}
		printOutput(*jsonOutput, tunnels)
	case "inspect":
		if flags.NArg() < 1 {
			log.Fatal("tunnel ID required")
		}
		var tunnel TunnelDTO
		if err := apiClient.Do(context.Background(), http.MethodGet, "/api/v1/tunnels/"+flags.Arg(0), nil, &tunnel); err != nil {
			log.Fatalf("inspect tunnel: %v", err)
		}
		printOutput(*jsonOutput, tunnel)
	case "create", "start":
		if flags.NArg() < 1 {
			log.Fatal("tunnel name/protocol required")
		}
		var tunnel TunnelDTO
		payload := map[string]string{"name": flags.Arg(0), "protocol": "http"}
		if flags.NArg() > 1 {
			payload["protocol"] = flags.Arg(1)
		}
		if err := apiClient.Do(context.Background(), http.MethodPost, "/api/v1/tunnels", payload, &tunnel); err != nil {
			log.Fatalf("create tunnel: %v", err)
		}
		printOutput(*jsonOutput, tunnel)
	case "stop", "revoke":
		if flags.NArg() < 1 {
			log.Fatal("tunnel ID required")
		}
		if err := apiClient.Do(context.Background(), http.MethodDelete, "/api/v1/tunnels/"+flags.Arg(0), nil, nil); err != nil {
			log.Fatalf("stop tunnel: %v", err)
		}
		fmt.Println("tunnel stopped")
	default:
		log.Fatalf("unknown action %q", action)
	}
}
