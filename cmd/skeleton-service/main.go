package main

import (
	"flag"
	"fmt"
	"skeleton/cmd/config"
)

const DEFAULT_CONFIG_PATH = "./config/config.yaml"

var (
	Version   = "undefined"
	Commit    = "undefined"
	Timestamp = "undefined"
)

func main() {
	configFile := flag.String("config", DEFAULT_CONFIG_PATH, "config file")
	flag.Parse()

	specs, err := config.LoadConfig(*configFile)
	if err != nil {
		panic(fmt.Errorf("Error loading config: %w", err))
	}
	specs.Version = Version
	specs.Commit = Commit
	specs.Timestamp = Timestamp
	fmt.Printf("specs: %+v", specs)
}
