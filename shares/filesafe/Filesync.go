package filesafe

import "os"

// 强制刷盘
func Fsync(fpath string) error {
	file, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer file.Close()
	return file.Sync()
}
