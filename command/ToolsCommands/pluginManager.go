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
	manifest, err := pluginsloader.LoadPluginManifest(pluginDir) // 获取插件清单
	if err != nil {
		return fmt.Errorf("加载插件清单失败: %w", err)
	}

	if strings.TrimSpace(manifest.Name) == "" { // 获取插件的名称
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

func listAllPlugins() ([]string, error) {
	cfg, err := config.LoadMisakaConfigure(defaultConfigPath)
	if err != nil {
		return nil, err
	}
	installed := cfg.Private.Storage.Plugins // 获取安装的所有插件
	return installed, nil
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

	imports, err := resolveEnabledPluginImports(".", cfg.Private.Storage.Plugins)
	if err != nil {
		return err
	}

	clilog.Info(fmt.Sprintf("插件列表: %v", imports))
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
		// 没有插件的情况 不需要使用 pluginbridge 了
		return nil
	}

	var builder strings.Builder

	builder.WriteString("package pluginbridge\n\n")
	builder.WriteString("// Code generated by misaka-tools. DO NOT EDIT.\n")
	builder.WriteString("import (\n")
	builder.WriteString("\tpluginsloader \"misakadb/plugins/pluginsLoader\"\n")
	builder.WriteString("\"misakadb/clilog\"\n")

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
