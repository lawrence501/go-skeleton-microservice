package config

import (
	"fmt"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

func LoadConfig(configFile string) (Specs, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
		return Specs{}, err
	}

	var config Specs
	if err := k.Unmarshal("spec", &config); err != nil {
		return Specs{}, err
	}
	fmt.Println("serviceName:", config.ServiceName)
	return config, nil
}
