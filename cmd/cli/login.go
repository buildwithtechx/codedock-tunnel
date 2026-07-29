package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"codedock.run/codedock-tunnel/internal/config"
)

func runLogin(cfg config.CLIConfig) {
	flags := flag.NewFlagSet("login", flag.ExitOnError)
	token := flags.String("agent-token", "", "agent token issued by the dashboard")
	apiKey := flags.String("api-key", cfg.APIKey, "API key for management commands")
	_ = flags.Parse(os.Args[2:])
	if *token == "" {
		log.Fatal("login requires --agent-token issued by the dashboard")
	}
	cfg.AgentToken = *token
	cfg.APIKey = *apiKey
	if err := config.SaveCLI(cfg); err != nil {
		log.Fatalf("save credentials: %v", err)
	}
	fmt.Println("credentials saved")
}
