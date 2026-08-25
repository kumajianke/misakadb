package judgment

import (
	"errors"
	"io/fs"
	"misakadb/clilog"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func fileExists(rootDir, targetName string) (string, bool, error) {
	var foundPath string

	if rootDir == targetName {
		return "", true, nil
	}

	err := filepath.WalkDir(rootDir, func(currentPath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() && entry.Name() == targetName {
			foundPath = currentPath
			return filepath.SkipAll
		}

		return nil
	})
	if err != nil {
		return "", false, err
	}

	return foundPath, foundPath != "", nil
}

func removeFileClear(folder string) error {
	all_file, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, file := range all_file {
		file_detail, err := file.Info()
		if err != nil {
			clilog.Error("[error] can not get target detail of file：" + err.Error())
			continue
		}
		filename := file_detail.Name()
		is_removed_file := strings.HasPrefix(filename, "MisakaRemove.")

		if is_removed_file {
			var err error
			if file_detail.IsDir() {
				err = os.RemoveAll(path.Join(folder, filename))
			} else {
				err = os.Remove(path.Join(folder, filename))

			}
			if err != nil { // 首次删除失败
				// 报错
				clilog.Error("[error] can not get target detail of file：" + err.Error())
				return err
			}

		} else {
			// 不是要删除的文件那就ISDIR递归进去
			if file_detail.IsDir() {
				if err := removeFileClear(path.Join(folder, filename)); err != nil {
					return err
				}
				continue
			}

		}
	}
	return nil
}

var AllRemoveCheck []string = []string{"./db-datas"}

func AppendAndGetRemoveCheck(new_path string) ([]string, error) {
	for _, p := range AllRemoveCheck {
		_, flag, err := fileExists(p, new_path)
		if err != nil {
			return nil, err
		} else if flag {
			return nil, errors.New("the path is in the allremovecheck or it is item children path!")
		}
	}

	if _, e := os.ReadDir(new_path); e != nil {
		return nil, e
	} else {
		AllRemoveCheck = append(AllRemoveCheck, new_path)
		return AllRemoveCheck, nil
	}
}

// 删除已经被 标记删除的文件
func NewThreadStartJudgmentClear() {
	for _, p := range AllRemoveCheck {
		go removeFileClear(p)
	}
}
