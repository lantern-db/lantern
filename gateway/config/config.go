package config

import (
	"errors"
	"os"
)

type GatewayConfig struct {
	LanternHost string
	LanternPort string
	GatewayPort string
}

func LoadGatewayConfig() (*GatewayConfig, error) {
	lanternHost := os.Getenv("LANTERN_HOST")
	if len(lanternHost) == 0 {
		return nil, errors.New("environment LANTERN_HOST should be set")
	}

	lanternPort := os.Getenv("LANTERN_PORT")
	if len(lanternPort) == 0 {
		return nil, errors.New("environment LANTERN_PORT should be set")
	}

	gatewayPort := os.Getenv("GATEWAY_PORT")
	if len(gatewayPort) == 0 {
		return nil, errors.New("environment GATEWAY_PORT")
	}

	return &GatewayConfig{
		LanternHost: lanternHost,
		LanternPort: lanternPort,
		GatewayPort: gatewayPort,
	}, nil
}
