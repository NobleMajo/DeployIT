package stringfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveFile_missingPath(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	p := filepath.Join(dir, "nope-not-created")
	if err := RemoveFile(p); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveFile_regularFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	p := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := RemoveFile(p); err != nil {
		t.Fatal(err)
	}
	if Exists(p) {
		t.Fatal("file should be gone")
	}
}

func TestRemoveFile_directory(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	sub := filepath.Join(dir, "d")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	inner := filepath.Join(sub, "f")
	if err := os.WriteFile(inner, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := RemoveFile(sub); err != nil {
		t.Fatal(err)
	}
	if Exists(sub) {
		t.Fatal("directory should be removed")
	}
}
