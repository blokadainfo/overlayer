package main

import (
	"log"

	"github.com/blokadainfo/overlayer/config"
	"github.com/blokadainfo/overlayer/stream"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config from environment: %v", err)
	}

	for _, sc := range c.Streams {
		go stream.StartStream(c.Overlay, sc) // Start stream in a new goroutine
	}

	select {} // Block forever, so all goroutines can run
}
