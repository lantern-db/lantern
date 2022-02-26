package config

import (
	"errors"
	"os"
)

type PrometheusConfig struct {
	Port string
}

func LoadPrometheusConfig() (*PrometheusConfig, error) {
	prometheusPort := os.Getenv("PROMETHEUS_PORT")
	if len(prometheusPort) == 0 {
		return nil, errors.New("environment PROMETHEUS_PORT should be set")
	}

	return &PrometheusConfig{
		Port: prometheusPort,
	}, nil
}
