package main

import (
	"log"

	"codedock.run/codedock-tunnel/internal/config"
)

var version = "dev"

func main() {
	if _, err := config.LoadCron(); err != nil {
		log.Fatal(err)
	}
}
