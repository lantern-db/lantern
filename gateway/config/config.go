package config

import (
	"errors"
	"os"
)

type GatewayConfig struct {
	LanternHost   string
	LanternPort   string
	GatewayPort   string
	AllowedOrigin string
}

func LoadGatewayConfig() (*GatewayConfig, error) {
	lanternHost := os.Getenv("LANTERN_HOST")
	if len(lanternHost) == 0 {
		return nil, errors.New("environment LANTERN_HOST should be set")
	}

	lanternPort := os.Getenv("LANTERN_PORT")
	if len(lanternPort) == 0 {
		return nil, errors.New("environment LANTERN_PORT must be set")
	}

	gatewayPort := os.Getenv("GATEWAY_PORT")
	if len(gatewayPort) == 0 {
		return nil, errors.New("environment GATEWAY_PORT must be set")
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if len(allowedOrigin) == 0 {
		return nil, errors.New("environment ALLOWED_ORIGIN must be set")
	}

	return &GatewayConfig{
		LanternHost:   lanternHost,
		LanternPort:   lanternPort,
		GatewayPort:   gatewayPort,
		AllowedOrigin: allowedOrigin,
	}, nil
}
