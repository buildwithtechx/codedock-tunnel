package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"codedock.run/codedock-tunnel/internal/config"
	"codedock.run/codedock-tunnel/pkg/client"
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
		log.Fatalf("unknown command %q", os.Args[1])
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

func printUsage() {
	fmt.Println("codedock-tunnel <command>")
	fmt.Println("commands: health, version")
}
