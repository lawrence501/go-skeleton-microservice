package config

import (
	"fmt"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func CreateLogger(lc LogConfig) (logger *zap.Logger, err error) {
	level, err := zap.ParseAtomicLevel(lc.Level)
	if err != nil {
		return nil, err
	}
	conf := zap.NewDevelopmentConfig()
	conf.Level = level
	conf.DisableStacktrace = !lc.Stacktrace
	conf.OutputPaths = []string{"stdout"}
	conf.EncoderConfig.LevelKey = "severity"
	conf.EncoderConfig.TimeKey = "timestamp"
	conf.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	conf.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	return conf.Build()
}
