package pluginsloader

import (
	"fmt"
	"misakadb/config"
	"strings"
)

var builtinPluginRegisters = map[string]func() error{}

// RegisterBuiltinPlugin 注册一个内建插件的启动注册函数。
func RegisterBuiltinPlugin(pluginName string, register func() error) {
        pluginName = normalizePluginName(pluginName)
	if pluginName == "" || register == nil {
		return
	}
	builtinPluginRegisters[pluginName] = register
}

// BootstrapPlugins 根据配置中启用的插件列表执行对应的 Register 函数。
func BootstrapPlugins() error {
	cfg := config.GetGlobalMisakaConfigure()
	if cfg == nil {
		return fmt.Errorf("global config is nil")
	}

	for _, pluginName := range cfg.Private.Storage.Plugins {
                normalizedPluginName := normalizePluginName(pluginName)
		if normalizedPluginName == "" {
			continue
		}

		register, ok := builtinPluginRegisters[normalizedPluginName]
		if !ok {
			continue
		}
		if err := register(); err != nil {
			return fmt.Errorf("bootstrap plugin %s failed: %w", strings.TrimSpace(pluginName), err)
		}
	}

	return nil
}
