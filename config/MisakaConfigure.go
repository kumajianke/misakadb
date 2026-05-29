package config

type MisakaConfigure struct {
	Network NetworkConfigure `yaml:"network" json:"network"`
	Service ServiceConfigure `yaml:"service" json:"service"`
	Private PrivateConfigure `yaml:"private" json:"private"`
}

type NetworkConfigure struct {
	Port       int    `yaml:"port" json:"port"`
	Address    string `yaml:"address" json:"address"`
	MaxConn    int    `yaml:"max_conn" json:"max_conn"`
	RetryCount int    `yaml:"retry_count" json:"retry_count"`
	RetryDelay int    `yaml:"retry_delay" json:"retry_delay"`
}

type ServiceConfigure struct {
	Version string `yaml:"version" json:"version"`
}

type PrivateConfigure struct {
	Storage StorageConfigure `yaml:"storage" json:"storage"`
}

type StorageConfigure struct {
	Path string `yaml:"path" json:"path"`
}
