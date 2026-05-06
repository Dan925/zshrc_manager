package backup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const layout = "2006-01-02_15-04-05"

type Backup struct {
	Path      string
	Timestamp time.Time
	SizeBytes int64
}

func backupDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home dir: %w", err)
	}
	return filepath.Join(home, ".zshrc.bak"), nil
}

func Create(src string) (string, error) {
	dir, err := backupDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating backup dir: %w", err)
	}
	dst := filepath.Join(dir, time.Now().Format(layout)+".zshrc")
	if err := copyFile(src, dst); err != nil {
		return "", fmt.Errorf("creating backup: %w", err)
	}
	return dst, nil
}

func List() ([]Backup, error) {
	dir, err := backupDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading backup dir: %w", err)
	}

	var backups []Backup
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".zshrc") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".zshrc")
		ts, err := time.ParseInLocation(layout, name, time.Local)
		if err != nil {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		backups = append(backups, Backup{
			Path:      filepath.Join(dir, e.Name()),
			Timestamp: ts,
			SizeBytes: info.Size(),
		})
	}
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.After(backups[j].Timestamp)
	})
	return backups, nil
}

func Restore(backupPath, dst string) error {
	if err := copyFile(backupPath, dst); err != nil {
		return fmt.Errorf("restoring backup: %w", err)
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
