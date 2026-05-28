package config

import (
	"encoding/json"
	"os"
	"sync"

	"misakadb/misaka_network"

	"gopkg.in/yaml.v3"
)

var (
	globalConfigure     *MisakaConfigure
	globalConfigureErr  error
	globalConfigureOnce sync.Once
)

func LoadMisakaConfigure(path string) (*MisakaConfigure, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &MisakaConfigure{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	applyDefaults(cfg)
	return cfg, nil
}

func ConvertServiceInfo(cfg *MisakaConfigure) *misaka_network.ServiceInfo {
	if cfg == nil {
		return nil
	}

	applyDefaults(cfg)
	return &misaka_network.ServiceInfo{
		Port:    cfg.Network.Port,
		Address: cfg.Network.Address,
	}
}

func LoadServiceInfo(path string) (*misaka_network.ServiceInfo, error) {
	cfg, err := LoadMisakaConfigure(path)
	if err != nil {
		return nil, err
	}
	return ConvertServiceInfo(cfg), nil
}

func InitGlobalMisakaConfigure(path string) (*MisakaConfigure, error) {
	globalConfigureOnce.Do(func() {
		globalConfigure, globalConfigureErr = LoadMisakaConfigure(path)
	})
	return globalConfigure, globalConfigureErr
}

func GetGlobalMisakaConfigure() *MisakaConfigure {
	return globalConfigure
}

func GetGlobalNetworkConfigure() *NetworkConfigure {
	if globalConfigure == nil {
		return nil
	}
	return &globalConfigure.Network
}

func GetGlobalServiceConfigure() *ServiceConfigure {
	if globalConfigure == nil {
		return nil
	}
	return &globalConfigure.Service
}

func GetGlobalServiceInfo() *misaka_network.ServiceInfo {
	return ConvertServiceInfo(globalConfigure)
}

func applyDefaults(cfg *MisakaConfigure) {
	if cfg.Network.Address == "" {
		cfg.Network.Address = "0.0.0.0"
	}
	if cfg.Network.Port == 0 {
		cfg.Network.Port = 10032
	}
	if cfg.Service.AllowCommand == nil {
		cfg.Service.AllowCommand = []string{}
	}
}

func ConvertConfigureToJSON(cfg *MisakaConfigure) string {
	if cfg == nil {
		return ""
	}

	applyDefaults(cfg)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}
