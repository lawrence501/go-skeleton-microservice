package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"skeleton/cmd/config"
	"skeleton/pkg/httpclient"
	"skeleton/pkg/opencensus"
	"time"

	"go.uber.org/zap"
)

const DEFAULT_CONFIG_PATH = "./config/config.yaml"

var (
	Version   = "undefined"
	Commit    = "undefined"
	Timestamp = "undefined"
)

func main() {
	// Config
	specs, err := loadConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading config: %w", err))
	}

	// Logging
	logger, err := config.CreateLogger(specs.Log)
	if err != nil {
		panic(fmt.Errorf("Error creating logger: %w", err))
	}
	logger = logger.With(zap.String("serviceName", specs.ServiceName), zap.String("version", specs.Version))
	zap.ReplaceGlobals(logger)
	ctx := context.WithValue(context.Background(), config.LoggingContextKey{}, logger)

	logger.Info("Starting service...")
	specJSON, err := json.Marshal(specs)
	if err != nil {
		logger.Panic("8Error while marshalling config JSON", zap.Error(err))
	}
	logger.Info("Config", zap.String("config", string(specJSON)))

	// HTTP/opencensus config
	setupDefaultTransport()
	createHTTPClient()

	// Healthcheck endpoint setup
	healthHandler := createHealthHandler(ctx, specs)
}

func createHealthHandler(ctx context.Context, specs config.Specs) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := json.Marshal(struct {
			ServiceName string `json:"serviceName"`
			Version     string `json:"version"`
			Status      string `json:"status"`
		}{
			ServiceName: specs.ServiceName,
			Version:     specs.Version,
			Status:      "UP",
		})
		if err != nil {
			zap.L().Error("Error serving healthcheck", zap.Error(err))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
}

func createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * time.Duration(60),
		Transport: opencensus.WrapTransport(
			httpclient.LoggingRoundTrip{
				Base: http.DefaultTransport,
			}),
	}
}

func setupDefaultTransport() {
	transport := http.DefaultTransport.(*http.Transport)
	transport.MaxIdleConns = 100
	transport.MaxIdleConnsPerHost = 100
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
