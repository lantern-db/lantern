package model

import (
	"log"
	"os"
	"strconv"
	"time"
)

type LanternServerConfig struct {
	Port          string
	Ttl           time.Duration
	FlushInterval time.Duration
}

func LoadServerConfig() *LanternServerConfig {
	flushInterval, err := strconv.Atoi(os.Getenv("LANTERN_FLUSH_INTERVAL"))
	if err != nil {
		log.Fatalf("flush interval parse failed: %v", err)
	}

	lanternPort := os.Getenv("LANTERN_PORT")

	ttl, err := strconv.Atoi(os.Getenv("LANTERN_TTL"))
	if err != nil {
		log.Fatalf("ttl parse error: %v", err)
	}

	return &LanternServerConfig{
		Port:          lanternPort,
		Ttl:           time.Duration(ttl) * time.Second,
		FlushInterval: time.Duration(flushInterval) * time.Second,
	}
}
