package backup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateAndList(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	src := filepath.Join(tmpHome, ".zshrc")
	os.WriteFile(src, []byte("alias gs='git status'\n"), 0644)

	backupPath, err := Create(src)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if _, err := os.Stat(backupPath); err != nil {
		t.Errorf("backup file not found: %v", err)
	}

	backups, err := List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(backups) != 1 {
		t.Fatalf("expected 1 backup, got %d", len(backups))
	}
}

func TestRestore(t *testing.T) {
	origHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	src := filepath.Join(tmpHome, ".zshrc")
	os.WriteFile(src, []byte("alias gs='git status'\n"), 0644)

	backupPath, _ := Create(src)
	os.WriteFile(src, []byte("alias gs='git diff'\n"), 0644)

	if err := Restore(backupPath, src); err != nil {
		t.Fatalf("Restore() error: %v", err)
	}

	content, _ := os.ReadFile(src)
	if string(content) != "alias gs='git status'\n" {
		t.Errorf("restored content = %q, want %q", content, "alias gs='git status'\n")
	}
}
