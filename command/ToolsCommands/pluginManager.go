package toolscommands

import (
	"fmt"
	"io/fs"
	"misakadb/clilog"
	"misakadb/config"
	pluginsloader "misakadb/plugins/pluginsLoader"
	"misakadb/shares"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "./profiles/misaka.yaml"
const defaultBuildScriptPath = "./build.sh"
const defaultBuildScriptWindowsPath = "./build.bat"
const defaultBridgeDirPath = "./plugins/pluginbridge"
const defaultBridgeFilePath = "./plugins/pluginbridge/bridge_gen.go"
const defaultModulePath = "misakadb"

func installPlugin(pluginDir string) error {
	cfg, err := config.LoadMisakaConfigure(defaultConfigPath)
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

	pluginName, configuredEntry, err := buildConfiguredEntryFromPluginDir(".", pluginDir)
	if err != nil {
		return err
	}
	for _, installed := range cfg.Private.Storage.Plugins {
		if pluginsloader.NormalizeConfiguredPluginName(installed) == pluginName {
			return nil
		}
	}

	cfg.Private.Storage.Plugins = append(cfg.Private.Storage.Plugins, configuredEntry)
	return saveMisakaConfigure(defaultConfigPath, cfg)
}

func uninstallPlugin(pluginName string) error {
	cfg, err := config.LoadMisakaConfigure(defaultConfigPath)
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

	pluginName = pluginsloader.NormalizeConfiguredPluginName(pluginName)
	filtered := make([]string, 0, len(cfg.Private.Storage.Plugins))
	for _, installed := range cfg.Private.Storage.Plugins {
		if pluginsloader.NormalizeConfiguredPluginName(installed) != pluginName {
			filtered = append(filtered, installed)
		}
	}
	cfg.Private.Storage.Plugins = filtered
	return saveMisakaConfigure(defaultConfigPath, cfg)
}

func listAllPlugins() ([]string, error) {
	cfg, err := config.LoadMisakaConfigure(defaultConfigPath)
	if err != nil {
		return nil, err
	}
	installed := make([]string, 0, len(cfg.Private.Storage.Plugins))
	for _, plugin := range cfg.Private.Storage.Plugins {
		pluginName := pluginsloader.NormalizeConfiguredPluginName(plugin)
		if pluginName == "" {
			continue
		}
		installed = append(installed, pluginName)
	}
	return installed, nil
}

func syncInstalledPlugins() error {
	cfg, err := config.LoadMisakaConfigure(defaultConfigPath)
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

	configuredPlugins := uniqueConfiguredPlugins(cfg.Private.Storage.Plugins)
	pluginDirs, err := resolveConfiguredPluginDirs(".", configuredPlugins)
	if err != nil {
		return err
	}

	refreshedPlugins := make([]string, 0, len(pluginDirs))
	for _, pluginDir := range pluginDirs {
		_, configuredEntry, err := buildConfiguredEntryFromPluginDir(".", pluginDir)
		if err != nil {
			return fmt.Errorf("重新安装插件 %s 失败: %w", pluginDir, err)
		}
		refreshedPlugins = append(refreshedPlugins, configuredEntry)
	}

	cfg.Private.Storage.Plugins = refreshedPlugins
	if err := saveMisakaConfigure(defaultConfigPath, cfg); err != nil {
		return fmt.Errorf("写入刷新后的插件配置失败: %w", err)
	}
	return nil
}

func rebuildAfterPluginChange() error {
	if err := syncPluginBridge(); err != nil {
		return fmt.Errorf("同步插件桥接文件失败: %w", err)
	}

	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("未检测到 Go 开发环境，请先安装 Go 并确保 `go` 已加入 PATH")
	}

	var cmd *exec.Cmd
	if shares.IsWindows() {
		scriptPath, err := filepath.Abs(defaultBuildScriptWindowsPath)
		if err != nil {
			return fmt.Errorf("解析 Windows 构建脚本路径失败: %w", err)
		}
		if _, err := os.Stat(scriptPath); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Windows 环境缺少构建脚本: %s", scriptPath)
			}
			return fmt.Errorf("检查 Windows 构建脚本失败: %w", err)
		}
		cmd = exec.Command("cmd", "/c", scriptPath)
	} else {
		scriptPath, err := filepath.Abs(defaultBuildScriptPath)
		if err != nil {
			return fmt.Errorf("解析构建脚本路径失败: %w", err)
		}
		cmd = exec.Command("sh", scriptPath)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type pluginImportEntry struct {
	Name         string
	ImportPath   string
	ImportAlias  string
	BootFunction string
}

func syncPluginBridge() error {
	cfg, err := config.LoadMisakaConfigure(defaultConfigPath)
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

	imports, err := resolveEnabledPluginImports(".", uniqueConfiguredPlugins(cfg.Private.Storage.Plugins))
	if err != nil {
		return err
	}

	clilog.Info(fmt.Sprintf("插件列表: %v", imports))
	return writePluginBridgeFile(imports)
}

func resolveConfiguredPluginDirs(repoRoot string, enabledPlugins []string) ([]string, error) {
	pluginDirs, err := discoverPluginDirs(repoRoot)
	if err != nil {
		return nil, err
	}

	dirByName := make(map[string]string, len(pluginDirs))
	for _, pluginDir := range pluginDirs {
		manifest, err := pluginsloader.LoadPluginManifest(pluginDir)
		if err != nil {
			return nil, fmt.Errorf("加载插件清单失败 %s: %w", pluginDir, err)
		}

		pluginName := strings.TrimSpace(manifest.Name)
		if pluginName == "" {
			return nil, fmt.Errorf("插件清单缺少 name: %s", pluginDir)
		}
		if existingDir, ok := dirByName[pluginName]; ok && existingDir != pluginDir {
			return nil, fmt.Errorf("发现重名插件 %q: %s 与 %s", pluginName, existingDir, pluginDir)
		}
		dirByName[pluginName] = pluginDir
	}

	resolvedDirs := make([]string, 0, len(enabledPlugins))
	for _, rawPlugin := range enabledPlugins {
		entry := pluginsloader.ParseConfiguredPluginEntry(rawPlugin)
		if entry.Name == "" {
			continue
		}

		if entry.DirHint != "" {
			pluginDir := filepath.Clean(filepath.Join(repoRoot, filepath.FromSlash(entry.DirHint)))
			manifest, err := pluginsloader.LoadPluginManifest(pluginDir)
			if err != nil {
				return nil, fmt.Errorf("已安装插件 %q 的目录定位 %q 无效: %w", entry.Name, entry.DirHint, err)
			}
			if strings.TrimSpace(manifest.Name) != entry.Name {
				return nil, fmt.Errorf("已安装插件 %q 的目录定位 %q 指向了插件 %q", entry.Name, entry.DirHint, strings.TrimSpace(manifest.Name))
			}
			resolvedDirs = append(resolvedDirs, pluginDir)
			continue
		}

		pluginDir, ok := dirByName[entry.Name]
		if !ok {
			clilog.Warning(fmt.Sprintf("跳过无法定位目录的旧插件配置: %s", entry.Name))
			continue
		}
		resolvedDirs = append(resolvedDirs, pluginDir)
	}

	return resolvedDirs, nil
}

func uniqueConfiguredPlugins(pluginEntries []string) []string {
	unique := make([]string, 0, len(pluginEntries))
	seen := make(map[string]struct{}, len(pluginEntries))
	for _, rawEntry := range pluginEntries {
		pluginName := pluginsloader.NormalizeConfiguredPluginName(rawEntry)
		if pluginName == "" {
			continue
		}
		if _, ok := seen[pluginName]; ok {
			continue
		}
		seen[pluginName] = struct{}{}
		unique = append(unique, rawEntry)
	}
	return unique
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

		bootImportPath, bootFunction, err := resolvePluginBootReference(repoRoot, pluginDir, manifest.Boot)
		if err != nil {
			return nil, fmt.Errorf("解析插件 %q 的 boot 配置失败: %w", pluginName, err)
		}

		entry := pluginImportEntry{
			Name:         pluginName,
			ImportPath:   bootImportPath,
			BootFunction: bootFunction,
		}

		if existed, ok := manifestByName[pluginName]; ok && existed.ImportPath != entry.ImportPath {
			return nil, fmt.Errorf("发现重名插件 %q: %s 与 %s", pluginName, existed.ImportPath, entry.ImportPath)
		}
		manifestByName[pluginName] = entry
	}

	resolvedImports := make([]pluginImportEntry, 0, len(enabledPlugins))
	for _, rawPlugin := range enabledPlugins {
		entryConfig := pluginsloader.ParseConfiguredPluginEntry(rawPlugin)
		if entryConfig.Name == "" {
			continue
		}

		if entryConfig.DirHint != "" {
			pluginDir := filepath.Clean(filepath.Join(repoRoot, filepath.FromSlash(entryConfig.DirHint)))
			manifest, err := pluginsloader.LoadPluginManifest(pluginDir)
			if err != nil {
				return nil, fmt.Errorf("已启用插件 %q 的目录定位 %q 无效: %w", entryConfig.Name, entryConfig.DirHint, err)
			}
			if strings.TrimSpace(manifest.Name) != entryConfig.Name {
				return nil, fmt.Errorf("已启用插件 %q 的目录定位 %q 指向了插件 %q", entryConfig.Name, entryConfig.DirHint, strings.TrimSpace(manifest.Name))
			}

			bootImportPath, bootFunction, err := resolvePluginBootReference(repoRoot, pluginDir, manifest.Boot)
			if err != nil {
				return nil, fmt.Errorf("解析插件 %q 的 boot 配置失败: %w", entryConfig.Name, err)
			}
			resolvedImports = append(resolvedImports, pluginImportEntry{
				Name:         entryConfig.Name,
				ImportPath:   bootImportPath,
				BootFunction: bootFunction,
			})
			continue
		}

		entry, ok := manifestByName[entryConfig.Name]
		if !ok {
			return nil, fmt.Errorf("已启用插件 %q 未找到对应的 plugin.yaml，无法生成桥接导入", entryConfig.Name)
		}
		resolvedImports = append(resolvedImports, entry)
	}

	sort.Slice(resolvedImports, func(i, j int) bool {
		return resolvedImports[i].ImportPath < resolvedImports[j].ImportPath
	})

	for i := range resolvedImports {
		resolvedImports[i].ImportAlias = fmt.Sprintf("plugin_%d", i)
	}

	return resolvedImports, nil
}

func resolvePluginBootReference(repoRoot string, pluginDir string, boot string) (string, string, error) {
	boot = strings.TrimSpace(boot)
	if boot == "" {
		return "", "", fmt.Errorf("boot 不能为空")
	}
	if !strings.HasPrefix(boot, "./") {
		return "", "", fmt.Errorf("boot 必须以 ./ 开头: %s", boot)
	}

	separator := strings.LastIndex(boot, "/")
	if separator <= 1 || separator == len(boot)-1 {
		return "", "", fmt.Errorf("boot 格式必须为 ./path/to/file.go/Func(): %s", boot)
	}

	fileRef := boot[:separator]
	functionRef := boot[separator+1:]
	if !strings.HasSuffix(functionRef, "()") {
		return "", "", fmt.Errorf("boot 函数必须以 () 结尾: %s", boot)
	}
	functionName := strings.TrimSuffix(functionRef, "()")
	if functionName == "" {
		return "", "", fmt.Errorf("boot 函数名不能为空: %s", boot)
	}

	filePath := filepath.Clean(filepath.Join(pluginDir, fileRef))
	relToPlugin, err := filepath.Rel(pluginDir, filePath)
	if err != nil {
		return "", "", fmt.Errorf("计算 boot 相对路径失败: %w", err)
	}
	if strings.HasPrefix(relToPlugin, "..") {
		return "", "", fmt.Errorf("boot 不能跳出插件目录: %s", boot)
	}
	if filepath.Ext(filePath) != ".go" {
		return "", "", fmt.Errorf("boot 必须指向 .go 文件: %s", boot)
	}

	importDir := filepath.Dir(filePath)
	relImportDir, err := filepath.Rel(repoRoot, importDir)
	if err != nil {
		return "", "", fmt.Errorf("计算 boot 导入路径失败: %w", err)
	}
	if strings.HasPrefix(relImportDir, "..") {
		return "", "", fmt.Errorf("boot 文件必须位于项目根目录内: %s", boot)
	}

	importPath := defaultModulePath
	if relImportDir != "." {
		importPath += "/" + filepath.ToSlash(relImportDir)
	}

	return importPath, functionName, nil
}

func buildPluginDirHint(repoRoot string, pluginDir string) (string, error) {
	relPath, err := filepath.Rel(repoRoot, pluginDir)
	if err != nil {
		return "", fmt.Errorf("计算插件目录相对路径失败 %s: %w", pluginDir, err)
	}
	if strings.HasPrefix(relPath, "..") {
		return "", fmt.Errorf("插件目录必须位于项目根目录内: %s", pluginDir)
	}
	return filepath.ToSlash(filepath.Clean(relPath)), nil
}

func buildConfiguredEntryFromPluginDir(repoRoot string, pluginDir string) (string, string, error) {
	manifest, err := pluginsloader.LoadPluginManifest(pluginDir)
	if err != nil {
		return "", "", fmt.Errorf("加载插件清单失败: %w", err)
	}

	pluginName := strings.TrimSpace(manifest.Name)
	if pluginName == "" {
		return "", "", fmt.Errorf("插件清单缺少 name: %s", pluginDir)
	}

	pluginDirHint, err := buildPluginDirHint(repoRoot, pluginDir)
	if err != nil {
		return "", "", err
	}
	return pluginName, pluginsloader.BuildConfiguredPluginEntry(pluginName, pluginDirHint), nil
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

/**
* 生成代码文件，以生成对应插件的运行入口
 */
func writePluginBridgeFile(imports []pluginImportEntry) error {
	if err := os.MkdirAll(defaultBridgeDirPath, 0700); err != nil {
		return err
	}

	if len(imports) == 0 {
		if err := os.Remove(defaultBridgeFilePath); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	var builder strings.Builder

	builder.WriteString("package pluginbridge\n\n")
	builder.WriteString("// Code generated by misaka-tools. DO NOT EDIT.\n")
	builder.WriteString("import (\n")
	builder.WriteString("\tpluginsloader \"misakadb/plugins/pluginsLoader\"\n")
	builder.WriteString("\t\"misakadb/clilog\"\n")

	for _, entry := range imports {
		builder.WriteString(fmt.Sprintf("\t%s %q\n", entry.ImportAlias, entry.ImportPath))
	}
	builder.WriteString(")\n\n")
	builder.WriteString("func RegisterBuiltinPlugins() {\n")
	builder.WriteString("\tclilog.Info(\"Load The All Plugins...\")\n")

	for _, entry := range imports {
		builder.WriteString(fmt.Sprintf("\tpluginsloader.RegisterBuiltinPlugin(%q, %s.%s)\n", entry.Name, entry.ImportAlias, entry.BootFunction))
	}
	builder.WriteString("}\n")

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
