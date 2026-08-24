package filesafe

import (
	generashares "misakadb/shares"
	"os"
)

// 强制刷盘
func Fsync(fpath string) error {
	// Windows does not support fsync on directory handles.
	if generashares.IsWindows() {
		return nil
	}
	file, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer file.Close()
	return file.Sync()
}
