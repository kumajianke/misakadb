package config

type MisakaConfigure struct {
	Network NetworkConfigure `yaml:"network"`
	Service ServiceConfigure `yaml:"service"`
}

type NetworkConfigure struct {
	Port       int    `yaml:"port"`
	Address    string `yaml:"address"`
	MaxConn    int    `yaml:"max_conn"`
	RetryCount int    `yaml:"retry_count"`
	RetryDelay int    `yaml:"retry_delay"`
}

type ServiceConfigure struct {
	Version      string   `yaml:"version"`
	AllowCommand []string `yaml:"allow_command"`
}
