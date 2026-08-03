package global_lock

import (
	"path/filepath"
	"strings"
)

const fileHandleLockPrefix = "rw-filehandle-"

// GetFileHandleLockName 返回统一的文件句柄锁名。
func GetFileHandleLockName(handleFilePath string) string {
	handleFilePath = strings.TrimSpace(handleFilePath)
	if handleFilePath == "" {
		return fileHandleLockPrefix
	}

	absPath, err := filepath.Abs(handleFilePath)
	if err != nil {
		return fileHandleLockPrefix + filepath.Clean(handleFilePath)
	}
	return fileHandleLockPrefix + filepath.Clean(absPath)
}

// LockFileHandle 获取某个文件句柄对应的全局写锁。
func LockFileHandle(handleFilePath string) (func(), error) {
	_, unlock, err := GetOrStoreGlobalLock(GetFileHandleLockName(handleFilePath), "lock")
	return unlock, err
}
