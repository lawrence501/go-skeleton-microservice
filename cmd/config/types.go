package config

type Specs struct {
	ServiceName string         `koanf:"serviceName" json:"serviceName"`
	Version     string         `koanf:"version" json:"version"`
	Commit      string         `koanf:"commit" json:"commit"`
	Timestamp   string         `koanf:"timestamp" json:"timestamp"`
	Host        string         `koanf:"host" json:"host"`
	Port        int            `koanf:"port" json:"port"`
	DB          PostgresConfig `koanf:"db" json:"db"`
	Log         LogConfig      `koanf:"log" json:"log"`
}

type PostgresConfig struct {
	Host        string       `koanf:"host" json:"host"`
	Port        string       `koanf:"port" json:"port"`
	DBName      string       `koanf:"dbName" json:"dbName"`
	User        string       `koanf:"user" json:"user"`
	Password    SecretString `koanf:"password" json:"password"`
	MaxConn     int          `koanf:"maxConn" json:"maxConn"`
	MaxIdleConn int          `koanf:"maxIdleConn" json:"maxIdleConn"`
}

type LogConfig struct {
	Level      string `koanf:"level" json:"level"`
	Stacktrace bool   `koanf:"stacktrace" json:"stacktrace"`
}

type LoggingContextKey struct{}
