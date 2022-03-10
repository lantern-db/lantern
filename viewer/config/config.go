package config

import (
	"log"
	"os"
)

type ViewerConfig struct {
	Endpoint string `json:"endpoint,omitempty"`
}

func Load() *ViewerConfig {
	endpoint := os.Getenv("LANTERN_ENDPOINT")
	if len(endpoint) == 0 {
		log.Panicln("environment LANTERN_ENDPOINT must be set")
	}
	return &ViewerConfig{endpoint}
}
