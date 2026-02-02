package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zip"
)

// CalcFileMD5 计算文件MD5值
func CalcFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Unzip 解压zip文件到目标目录
// src：zip文件路径
// dst：解压目标目录
func Unzip(src, dst string) error {
	// 确保目标目录存在
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		// 拼接解压后的文件路径
		filePath := filepath.Join(dst, file.Name)
		// 防止路径遍历攻击（如文件名为../xxx）
		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return os.ErrInvalid
		}

		if file.FileInfo().IsDir() {
			// 创建目录
			_ = os.MkdirAll(filePath, file.Mode())
			continue
		}

		// 创建文件所在目录
		if err := os.MkdirAll(filepath.Dir(filePath), file.Mode()); err != nil {
			return err
		}

		// 打开目标文件
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// 打开zip内的文件
		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// 复制文件内容
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetFileSize 获取文件大小（字节）
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
