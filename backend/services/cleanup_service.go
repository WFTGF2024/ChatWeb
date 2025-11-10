package services

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"
)

// CleanupService 最小实现：可扩展 Start/Stop 定时调用 CleanOldFiles
type CleanupService struct {
	stop chan struct{}
}

func NewCleanupService() *CleanupService {
	return &CleanupService{stop: make(chan struct{})}
}

func (s *CleanupService) Start(interval time.Duration, days int) {
	if interval <= 0 {
		interval = time.Hour
	}
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = s.CleanOldFiles(days)
			case <-s.stop:
				return
			}
		}
	}()
}

func (s *CleanupService) Run(ctx context.Context, interval time.Duration, days int) {
	if interval <= 0 {
		interval = time.Hour
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_ = s.CleanOldFiles(days)
		case <-s.stop:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (s *CleanupService) Stop() {
	select {
	case <-s.stop:
	default:
		close(s.stop)
	}
}

// CleanOldFiles 删除 storage/chat 下超过 days 天的文件，并清理空目录
func (s *CleanupService) CleanOldFiles(days int) error {
	storagePath := "./storage/chat"
	if _, err := os.Stat(storagePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -days)
	return filepath.Walk(storagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil {
			return nil
		}
		// 删除超过指定天数的文件
		if !info.IsDir() && info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err != nil {
				return err
			}
			return nil
		}
		// 删除空目录
		if info.IsDir() {
			empty, _ := isDirEmpty(path)
			if empty {
				return os.Remove(path)
			}
		}
		return nil
	})
}

func isDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
