package model

import (
	"os"
	"strconv"
	"time"
)

type LanternServerConfig struct {
	Port          string
	Ttl           time.Duration
	FlushInterval time.Duration
}

func LoadServerConfig() (*LanternServerConfig, error) {
	flushInterval, err := strconv.Atoi(os.Getenv("LANTERN_FLUSH_INTERVAL"))
	if err != nil {
		return nil, err
	}

	lanternPort := os.Getenv("LANTERN_PORT")
	ttl, err := strconv.Atoi(os.Getenv("LANTERN_TTL"))
	if err != nil {
		return nil, err
	}

	return &LanternServerConfig{
		Port:          lanternPort,
		Ttl:           time.Duration(ttl) * time.Second,
		FlushInterval: time.Duration(flushInterval) * time.Second,
	}, nil
}
