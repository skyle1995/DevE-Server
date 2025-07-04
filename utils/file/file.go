package file

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Exists 检查文件或目录是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDir 检查路径是否为目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile 检查路径是否为文件
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// CreateDir 创建目录，如果不存在
func CreateDir(path string) error {
	if !Exists(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// ReadFile 读取文件内容
func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// ReadFileAsString 读取文件内容为字符串
func ReadFileAsString(path string) (string, error) {
	data, err := ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadLines 读取文件内容为字符串切片，每行一个元素
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// WriteFile 写入内容到文件
func WriteFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0644)
}

// WriteFileAsString 写入字符串到文件
func WriteFileAsString(path string, content string) error {
	return WriteFile(path, []byte(content))
}

// AppendToFile 追加内容到文件
func AppendToFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// AppendStringToFile 追加字符串到文件
func AppendStringToFile(path string, content string) error {
	return AppendToFile(path, []byte(content))
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// MoveFile 移动文件
func MoveFile(src, dst string) error {
	err := CopyFile(src, dst)
	if err != nil {
		return err
	}
	return os.Remove(src)
}

// DeleteFile 删除文件
func DeleteFile(path string) error {
	return os.Remove(path)
}

// ListFiles 列出目录下的所有文件
func ListFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filePaths []string
	for _, file := range files {
		filePaths = append(filePaths, filepath.Join(dir, file.Name()))
	}

	return filePaths, nil
}

// ListFilesRecursive 递归列出目录下的所有文件
func ListFilesRecursive(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetFileExt 获取文件扩展名
func GetFileExt(path string) string {
	return filepath.Ext(path)
}

// GetFileName 获取文件名（不含扩展名）
func GetFileName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// GetFilePath 获取文件所在目录
func GetFilePath(path string) string {
	return filepath.Dir(path)
}

// ReadJSON 从文件中读取JSON数据
func ReadJSON(path string, v interface{}) error {
	data, err := ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// WriteJSON 将JSON数据写入文件
func WriteJSON(path string, v interface{}, indent bool) error {
	var data []byte
	var err error

	if indent {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		return err
	}

	return WriteFile(path, data)
}
