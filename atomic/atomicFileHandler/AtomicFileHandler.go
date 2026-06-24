package atomic_file_handler

import (
	"errors"
	"fmt"
	"misakadb/clilog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
)

const (
	// 分片大小阈值 1MB
	ChunkSizeThreshold = 1 * 1024 * 1024
	// 单个分片大小
	ChunkSize = 1 * 1024 * 1024
)

// 获取分片存储的目录
func getChunkDir(filename string) string {
	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	return filepath.Join(dir, base+".chunk")
}

func syncDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	return d.Sync()
}

func parseChunkID(chunkFileName string) (int, error) {
	if filepath.Ext(chunkFileName) != ".tmp" {
		return 0, fmt.Errorf("invalid chunk file extension: %s", chunkFileName)
	}

	idPart := chunkFileName[:len(chunkFileName)-len(".tmp")]
	if idPart == "" {
		return 0, fmt.Errorf("empty chunk file id: %s", chunkFileName)
	}

	id, err := strconv.Atoi(idPart)
	if err != nil {
		return 0, fmt.Errorf("invalid chunk file name: %s", chunkFileName)
	}
	if id < 0 {
		return 0, fmt.Errorf("negative chunk file id: %s", chunkFileName)
	}

	return id, nil
}

// 强制写入并同步文件
func AtomicSyncWriteFile(filename string, content []byte, perm os.FileMode) error {
	tempFile, err := os.CreateTemp(filepath.Dir(filename), ".tmp-*")
	// 创建一个临时文件
	if err != nil {
		return err
	}

	// 结束的时候 删除临时文件
	tempFilename := tempFile.Name()
	defer func() {
		tempFile.Close() // 关闭临时文件操作句柄
		if err != nil {  // 创建临时文件没有失败的话
			os.Remove(tempFilename) // 删除临时文件
		}
	}()
	if _, err = tempFile.Write(content); err != nil {
		return err
	}
	if err = tempFile.Sync(); err != nil {
		return err
	}

	// 原子操作 重命名临时文件到目标文件名
	if err = os.Rename(tempFilename, filename); err != nil {
		return err
	}

	// 强制落盘 目标文件所在目录 不做Rename可能丢失
	dir := filepath.Dir(filename)
	return syncDir(dir)
}

// 分片存储文件，返回分片数量
func ChunkAtomicSyncWriteFile(filename string, content []byte, perm os.FileMode) (int, error) {
	if len(content) == 0 {
		return 0, nil
	}

	// 小于阈值直接写入
	if len(content) < ChunkSizeThreshold {
		err := AtomicSyncWriteFile(filename, content, perm)
		return 0, err
	}

	// 创建分片目录
	chunkDir := getChunkDir(filename)
	workDir := chunkDir + ".tmp"

	if err := os.RemoveAll(workDir); err != nil {
		return 0, err
	}

	if err := os.MkdirAll(workDir, 0755); err != nil {
		return 0, err
	}

	// 计算分片数量
	totalChunks := (len(content) + ChunkSize - 1) / ChunkSize

	// 分片写入
	for i := 0; i < totalChunks; i++ {
		start := i * ChunkSize
		end := start + ChunkSize
		if end > len(content) {
			end = len(content)
		}

		chunkFilename := filepath.Join(workDir, fmt.Sprintf("%d.tmp", i))
		if err := AtomicSyncWriteFile(chunkFilename, content[start:end], perm); err != nil {
			return 0, err
		}
	}

	// 创建备份
	backupDir := chunkDir + ".bak"
	if err := os.RemoveAll(backupDir); err != nil { // 删除已有的备份
		return 0, err
	}

	hasExistingChunkDir := false // 是否存在同名文件

	// 将已有的文件进行备份 而不是直接删除 防止数据丢失
	if _, err := os.Stat(chunkDir); err == nil {
		hasExistingChunkDir = true
		if err := os.Rename(chunkDir, backupDir); err != nil {
			return 0, err
		}
	}

	// 尝试吧本次事务的数据持久化到结果来
	if err := os.Rename(workDir, chunkDir); err != nil {
		if hasExistingChunkDir {
			rollbackErr := os.Rename(backupDir, chunkDir)
			if rollbackErr != nil {
				clilog.Error("restore chunk backup failed:", rollbackErr)
				return 0, errors.Join(err, fmt.Errorf("restore chunk backup failed: %w", rollbackErr))
			}
		}
		return 0, err
	}

	// 清理备份(旧文件)
	if hasExistingChunkDir {
		if err := os.RemoveAll(backupDir); err != nil {
			return 0, err
		}
	}

	// 强制同步目录
	if err := syncDir(filepath.Dir(chunkDir)); err != nil {
		return 0, err
	}

	return totalChunks, nil
}

// 合并分片文件到目标文件并删除分片目录
func ChunkMergeFile(filename string, perm os.FileMode) error {
	chunkDir := getChunkDir(filename)

	// 检查分片目录是否存在
	if _, err := os.Stat(chunkDir); os.IsNotExist(err) {
		return fmt.Errorf("chunk directory not found: %s", chunkDir)
	}

	// 读取所有分片文件
	entries, err := os.ReadDir(chunkDir)
	if err != nil {
		return err
	}

	type chunkFileMeta struct {
		name string
		id   int
	}

	// 过滤并排序分片文件
	var chunkFiles []chunkFileMeta
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".tmp" {
			id, parseErr := parseChunkID(entry.Name())
			if parseErr != nil {
				return parseErr
			}
			chunkFiles = append(chunkFiles, chunkFileMeta{
				name: entry.Name(),
				id:   id,
			})
		}
	}

	// 按分片ID排序
	sort.Slice(chunkFiles, func(i, j int) bool {
		return chunkFiles[i].id < chunkFiles[j].id
	})

	if len(chunkFiles) == 0 {
		return fmt.Errorf("no chunk files found in %s", chunkDir)
	}

	// 创建临时目标文件
	tempFile, err := os.CreateTemp(filepath.Dir(filename), ".merge-tmp-*")
	if err != nil {
		return err
	}
	tempFilename := tempFile.Name()

	defer func() {
		tempFile.Close()
		if err != nil {
			os.Remove(tempFilename)
		}
	}()

	// 并发读取分片并按顺序写入
	type chunkData struct {
		index int
		data  []byte
		err   error
	}

	resultChan := make(chan chunkData, len(chunkFiles))
	var wg sync.WaitGroup

	// 并发读取所有分片
	for i, chunkFile := range chunkFiles {
		wg.Add(1)
		go func(idx int, file chunkFileMeta) {
			defer wg.Done()
			chunkPath := filepath.Join(chunkDir, file.name)
			data, readErr := os.ReadFile(chunkPath)
			resultChan <- chunkData{index: idx, data: data, err: readErr}
		}(i, chunkFile)
	}

	// 等待所有读取完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果并排序
	chunks := make([]chunkData, 0, len(chunkFiles))
	for chunk := range resultChan {
		if chunk.err != nil {
			err = chunk.err
			return err
		}
		chunks = append(chunks, chunk)
	}

	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].index < chunks[j].index
	})

	// 按顺序写入合并文件
	for _, chunk := range chunks {
		if _, err = tempFile.Write(chunk.data); err != nil {
			return err
		}
	}

	// 同步到磁盘
	if err = tempFile.Sync(); err != nil {
		return err
	}

	// 关闭临时文件
	if err = tempFile.Close(); err != nil {
		return err
	}

	// 原子重命名
	if err = os.Rename(tempFilename, filename); err != nil {
		return err
	}

	// 强制落盘目录
	dir := filepath.Dir(filename)
	if err = syncDir(dir); err != nil {
		return err
	}

	// 删除分片目录
	return os.RemoveAll(chunkDir)
}

// 读取分片文件内容
func ChunkRead(filename string) ([]byte, error) {
	if err := ChunkMergeFile(filename, 0700); err != nil {
		return nil, err
	}

	content, err := os.ReadFile(filename)

	defer func() {
		if err != nil {
			err = os.Remove(filename)
			if err != nil {
				// ∂_∂ 删除合并后的文件失败
				clilog.Error("删除合并后的文件失败:", err)
			}
		}
	}()

	return content, err
}
