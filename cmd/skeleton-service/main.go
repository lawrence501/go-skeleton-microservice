package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"skeleton/cmd/config"

	"go.uber.org/zap"
)

const DEFAULT_CONFIG_PATH = "./config/config.yaml"

var (
	Version   = "undefined"
	Commit    = "undefined"
	Timestamp = "undefined"
)

func main() {
	specs, err := loadConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading config: %w", err))
	}

	logger, err := config.CreateLogger(specs.Log)
	if err != nil {
		panic(fmt.Errorf("Error creating logger: %w", err))
	}
	logger = logger.With(zap.String("service", specs.ServiceName), zap.String("version", specs.Version))
	zap.ReplaceGlobals(logger)
	_ = context.WithValue(context.Background(), config.LoggingContextKey{}, logger) // replace with ctx

	logger.Info("Starting service...")
	specJSON, err := json.Marshal(specs)
	if err != nil {
		logger.Panic("Error while marshalling config JSON", zap.Error(err))
	}
	logger.Info("Config", zap.String("config", string(specJSON)))
}

func loadConfig() (config.Specs, error) {
	configFile := flag.String("config", DEFAULT_CONFIG_PATH, "config file")
	flag.Parse()
	specs, err := config.LoadConfig(*configFile)
	if err != nil {
		return config.Specs{}, err
	}
	specs.Version = Version
	specs.Commit = Commit
	specs.Timestamp = Timestamp
	return specs, nil
}
