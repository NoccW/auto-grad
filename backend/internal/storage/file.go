package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type FileStorage struct {
	uploadPath string
}

func NewFileStorage(uploadPath string) *FileStorage {
	return &FileStorage{
		uploadPath: uploadPath,
	}
}

func (fs *FileStorage) SaveFile(file *multipart.FileHeader, subfolder string) (string, error) {
	// 确保上传目录存在
	fullPath := filepath.Join(fs.uploadPath, subfolder)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fs.generateUniqueFilename() + ext

	// 保存文件
	dst := filepath.Join(fullPath, filename)
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// 返回相对路径
	return filepath.Join(subfolder, filename), nil
}

func (fs *FileStorage) GetFilePath(relativePath string) string {
	return filepath.Join(fs.uploadPath, relativePath)
}

func (fs *FileStorage) DeleteFile(relativePath string) error {
	fullPath := filepath.Join(fs.uploadPath, relativePath)
	return os.Remove(fullPath)
}

func (fs *FileStorage) generateUniqueFilename() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes) + "-" + time.Now().Format("20060102150405")
}
