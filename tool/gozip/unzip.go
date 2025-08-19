package gozip

import (
	"archive/zip"
	"fmt"
	"github.com/biu7/gokit/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(input, output string) ([]string, error) {
	archive, err := zip.OpenReader(input)
	if err != nil {
		panic(err)
	}
	defer func() {
		if closeErr := archive.Close(); closeErr != nil {
			log.Default.Error("[unzip] close gozip file error", "error", err)
		}
	}()

	// 创建解压目标目录
	if err = os.MkdirAll(output, 0755); err != nil {
		return nil, err
	}

	var files []string
	// 遍历 gozip 文件内的所有文件和目录
	for _, file := range archive.File {
		err = extractFile(output, file)
		if err != nil {
			return nil, err
		}
		files = append(files, file.Name)
	}
	return files, nil
}

func extractFile(output string, f *zip.File) error {
	path := filepath.Join(output, f.Name)
	if !strings.HasPrefix(path, filepath.Clean(output)+string(os.PathSeparator)) {
		return fmt.Errorf("%s: illegal file path", path)
	}
	if f.FileInfo().IsDir() {
		err := os.MkdirAll(path, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", path, err)
		}
		return nil
	}

	// 创建父目录
	if err := os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
		return err
	}

	// 创建文件
	dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dstFile.Close(); closeErr != nil {
			log.Default.Error("[unzip] close file error: ", "error", closeErr, "file", path)
		}
	}()

	// 读取文件内容
	file, err := f.Open()
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", f.Name, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Default.Error("[unzip] close file error: ", "error", closeErr, "file", f.Name)
		}
	}()
	_, err = io.Copy(dstFile, file)
	if err != nil {
		return err
	}
	return nil
}
