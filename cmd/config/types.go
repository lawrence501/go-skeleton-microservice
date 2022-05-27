package config

type Specs struct {
	ServiceName string         `koanf:"serviceName"`
	Version     string         `koanf:"version"`
	Commit      string         `koanf:"commit"`
	Timestamp   string         `koanf:"timestamp"`
	Host        string         `koanf:"host"`
	Port        int            `koanf:"port"`
	DB          PostgresConfig `koanf:"db"`
}

type PostgresConfig struct {
	Host        string       `koanf:"host"`
	Port        string       `koanf:"port"`
	DBName      string       `koanf:"dbName"`
	User        string       `koanf:"user"`
	Password    SecretString `koanf:"password,omitempty"`
	MaxConn     int          `koanf:"maxConn"`
	MaxIdleConn int          `koanf:"maxIdleConn"`
}
