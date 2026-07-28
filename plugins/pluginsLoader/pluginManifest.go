package pluginsloader

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const PluginManifestFileName = "plugin.yaml"

type PluginTaskAlias struct {
	Alias    string `yaml:"alias"`
	TaskType string `yaml:"task_type"`
}

type PluginManifest struct {
	Name     string            `yaml:"name"`
	Boot     string            `yaml:"boot"`
	Register string            `yaml:"register"`
	Aliases  []PluginTaskAlias `yaml:"aliases"`
}

/*
* 加载指定插件（通过路径）的清单文件（Manifest）并返回
 */
func LoadPluginManifest(pluginDir string) (*PluginManifest, error) {
	manifestPath := filepath.Join(pluginDir, PluginManifestFileName)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	manifest := &PluginManifest{}
	if err := yaml.Unmarshal(data, manifest); err != nil {
		return nil, err
	}
	return manifest, nil
}
