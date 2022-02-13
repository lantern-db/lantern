package model

import (
	"os"
	"strconv"
	"time"
)

type LanternServerConfig struct {
	Port          string
	FlushInterval time.Duration
}

func LoadServerConfig() (*LanternServerConfig, error) {
	flushInterval, err := strconv.Atoi(os.Getenv("LANTERN_FLUSH_INTERVAL"))
	if err != nil {
		return nil, err
	}

	lanternPort := os.Getenv("LANTERN_PORT")

	return &LanternServerConfig{
		Port:          lanternPort,
		FlushInterval: time.Duration(flushInterval) * time.Second,
	}, nil
}
