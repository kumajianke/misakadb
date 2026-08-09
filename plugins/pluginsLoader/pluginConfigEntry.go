package pluginsloader

import "strings"

const configuredPluginDirMarker = "|dir="

type ConfiguredPluginEntry struct {
	Raw     string
	Name    string
	DirHint string
}

func ParseConfiguredPluginEntry(raw string) ConfiguredPluginEntry {
	raw = strings.TrimSpace(raw)
	entry := ConfiguredPluginEntry{
		Raw:  raw,
		Name: raw,
	}
	if raw == "" {
		return entry
	}

	index := strings.LastIndex(raw, configuredPluginDirMarker)
	if index < 0 {
		return entry
	}

	name := strings.TrimSpace(raw[:index])
	dirHint := strings.TrimSpace(raw[index+len(configuredPluginDirMarker):])
	if name == "" || dirHint == "" {
		return entry
	}

	entry.Name = name
	entry.DirHint = dirHint
	return entry
}

func NormalizeConfiguredPluginName(raw string) string {
	return ParseConfiguredPluginEntry(raw).Name
}

func BuildConfiguredPluginEntry(name string, dirHint string) string {
	name = strings.TrimSpace(name)
	dirHint = strings.TrimSpace(dirHint)
	if name == "" {
		return ""
	}
	if dirHint == "" {
		return name
	}
	return name + " " + configuredPluginDirMarker + dirHint
}
