// backend/utils/storage.go
package utils

import (
    "bytes"
    "compress/gzip"
    "crypto/md5"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

var storagePath = "./storage/chat" // 存储根目录

func init() {
    // 确保存储目录存在
    if err := os.MkdirAll(storagePath, 0755); err != nil {
        panic(err)
    }
}

// 压缩并保存到本地文件
func SaveToLocal(content string) (string, error) {
    // 压缩内容
    var buf bytes.Buffer
    gz := gzip.NewWriter(&buf)
    if _, err := gz.Write([]byte(content)); err != nil {
        return "", err
    }
    if err := gz.Close(); err != nil {
        return "", err
    }

    // 生成文件名（使用时间戳和内容哈希）
    hash := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
    timestamp := time.Now().Format("20060102")
    filename := fmt.Sprintf("%s_%s.gz", timestamp, hash[:8])
    
    // 创建按日期分组的子目录
    dir := filepath.Join(storagePath, timestamp)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return "", err
    }

    // 保存文件
    filepath := filepath.Join(dir, filename)
    if err := os.WriteFile(filepath, buf.Bytes(), 0644); err != nil {
        return "", err
    }

    return filepath, nil
}

// 从本地文件读取并解压
func LoadFromLocal(filepath string) (string, error) {
    data, err := os.ReadFile(filepath)
    if err != nil {
        return "", err
    }

    gz, err := gzip.NewReader(bytes.NewReader(data))
    if err != nil {
        return "", err
    }
    defer gz.Close()

    content, err := io.ReadAll(gz)
    if err != nil {
        return "", err
    }

    return string(content), nil
}

// 删除本地文件
func DeleteLocalFile(filepath string) error {
    return os.Remove(filepath)
}
