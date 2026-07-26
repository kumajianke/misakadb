package toolscommands

import (
	"fmt"
	"io/fs"
	"misakadb/config"
	pluginsloader "misakadb/plugins/pluginsLoader"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "./profiles/misaka.yaml"
const defaultBuildScriptPath = "./build.sh"
const defaultBridgeDirPath = "./plugins/pluginbridge"
const defaultBridgeFilePath = "./plugins/pluginbridge/bridge_gen.go"
const defaultModulePath = "misakadb"

func installPlugin(pluginDir string) error {
	manifest, err := pluginsloader.LoadPluginManifest(pluginDir)
	if err != nil {
		return fmt.Errorf("加载插件清单失败: %w", err)
	}
	if strings.TrimSpace(manifest.Name) == "" {
		return fmt.Errorf("插件清单缺少 name: %s", pluginDir)
	}

	cfg, err := config.LoadMisakaConfigure(defaultConfigPath)
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

	pluginName := strings.TrimSpace(manifest.Name)
	for _, installed := range cfg.Private.Storage.Plugins {
		if strings.TrimSpace(installed) == pluginName {
			return nil
		}
	}

	cfg.Private.Storage.Plugins = append(cfg.Private.Storage.Plugins, pluginName)
	return saveMisakaConfigure(defaultConfigPath, cfg)
}

func uninstallPlugin(pluginName string) error {
	cfg, err := config.LoadMisakaConfigure(defaultConfigPath)
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

	pluginName = strings.TrimSpace(pluginName)
	filtered := make([]string, 0, len(cfg.Private.Storage.Plugins))
	for _, installed := range cfg.Private.Storage.Plugins {
		if strings.TrimSpace(installed) != pluginName {
			filtered = append(filtered, installed)
		}
	}
	cfg.Private.Storage.Plugins = filtered
	return saveMisakaConfigure(defaultConfigPath, cfg)
}

func rebuildAfterPluginChange() error {
	if err := syncPluginBridge(); err != nil {
		return fmt.Errorf("同步插件桥接文件失败: %w", err)
	}

	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("未检测到 Go 开发环境，请先安装 Go 并确保 `go` 已加入 PATH")
	}

	cmd := exec.Command("sh", defaultBuildScriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type pluginImportEntry struct {
	Name       string
	ImportPath string
}

func syncPluginBridge() error {
	cfg, err := config.LoadMisakaConfigure(defaultConfigPath)
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

	imports, err := resolveEnabledPluginImports(".", cfg.Private.Storage.Plugins)
	if err != nil {
		return err
	}

	return writePluginBridgeFile(imports)
}

func resolveEnabledPluginImports(repoRoot string, enabledPlugins []string) ([]pluginImportEntry, error) {
	pluginDirs, err := discoverPluginDirs(repoRoot)
	if err != nil {
		return nil, err
	}

	manifestByName := make(map[string]pluginImportEntry, len(pluginDirs))
	for _, pluginDir := range pluginDirs {
		manifest, err := pluginsloader.LoadPluginManifest(pluginDir)
		if err != nil {
			return nil, fmt.Errorf("加载插件清单失败 %s: %w", pluginDir, err)
		}

		pluginName := strings.TrimSpace(manifest.Name)
		if pluginName == "" {
			return nil, fmt.Errorf("插件清单缺少 name: %s", pluginDir)
		}

		relPath, err := filepath.Rel(repoRoot, pluginDir)
		if err != nil {
			return nil, fmt.Errorf("计算插件相对路径失败 %s: %w", pluginDir, err)
		}
		if strings.HasPrefix(relPath, "..") {
			return nil, fmt.Errorf("插件目录必须位于项目根目录内: %s", pluginDir)
		}

		entry := pluginImportEntry{
			Name:       pluginName,
			ImportPath: defaultModulePath + "/" + filepath.ToSlash(relPath),
		}

		if existed, ok := manifestByName[pluginName]; ok && existed.ImportPath != entry.ImportPath {
			return nil, fmt.Errorf("发现重名插件 %q: %s 与 %s", pluginName, existed.ImportPath, entry.ImportPath)
		}
		manifestByName[pluginName] = entry
	}

	resolvedImports := make([]pluginImportEntry, 0, len(enabledPlugins))
	for _, pluginName := range enabledPlugins {
		pluginName = strings.TrimSpace(pluginName)
		if pluginName == "" {
			continue
		}

		entry, ok := manifestByName[pluginName]
		if !ok {
			return nil, fmt.Errorf("已启用插件 %q 未找到对应的 plugin.yaml，无法生成桥接导入", pluginName)
		}
		resolvedImports = append(resolvedImports, entry)
	}

	sort.Slice(resolvedImports, func(i, j int) bool {
		return resolvedImports[i].ImportPath < resolvedImports[j].ImportPath
	})

	return resolvedImports, nil
}

func discoverPluginDirs(repoRoot string) ([]string, error) {
	dirs := make([]string, 0)
	err := filepath.WalkDir(repoRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			if path != repoRoot && shouldSkipPluginWalkDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Name() != pluginsloader.PluginManifestFileName {
			return nil
		}

		dirs = append(dirs, filepath.Dir(path))
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(dirs)
	return dirs, nil
}

func shouldSkipPluginWalkDir(name string) bool {
	if strings.HasPrefix(name, ".") {
		return true
	}

	switch name {
	case "bin", "dist", "node_modules":
		return true
	default:
		return false
	}
}

func writePluginBridgeFile(imports []pluginImportEntry) error {
	if err := os.MkdirAll(defaultBridgeDirPath, 0700); err != nil {
		return err
	}

	var builder strings.Builder
	builder.WriteString("package pluginbridge\n\n")
	builder.WriteString("// Code generated by misaka-tools. DO NOT EDIT.\n")

	if len(imports) == 0 {
		builder.WriteString("\n")
	} else {
		builder.WriteString("import (\n")
		for _, entry := range imports {
			builder.WriteString(fmt.Sprintf("\t_ %q // %s\n", entry.ImportPath, entry.Name))
		}
		builder.WriteString(")\n")
	}

	return os.WriteFile(defaultBridgeFilePath, []byte(builder.String()), 0600)
}

func saveMisakaConfigure(path string, cfg *config.MisakaConfigure) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
