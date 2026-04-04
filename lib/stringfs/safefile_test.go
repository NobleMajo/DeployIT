package stringfs

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestSafeWriteFileBytes_success(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	target := filepath.Join(dir, "out.txt")
	content := []byte("hello-safe-write")
	if err := SafeWriteFileBytes(target, content, 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Fatalf("got %q", got)
	}
	tmp := filepath.Join(dir, ".tmp_out.txt")
	if Exists(tmp) {
		t.Fatal("temp file should be renamed away")
	}
}

func TestSafeWriteFileBytes_writeFails(t *testing.T) {
	t.Parallel()
	target := filepath.Join(t.TempDir(), "nope", "missing", "out.txt")
	err := SafeWriteFileBytes(target, []byte("x"), 0o644)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRemoveTmpSafeFile_removesTempSibling(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	finalPath := filepath.Join(dir, "out.txt")
	tmpPath := filepath.Join(dir, ".tmp_out.txt")
	if err := os.WriteFile(tmpPath, []byte("tmp"), 0o600); err != nil {
		t.Fatal(err)
	}
	RemoveTmpSafeFile(finalPath)
	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Fatalf("tmp file should be removed: %v", err)
	}
}

func TestRemoveTmpSafeFile_noTempOk(t *testing.T) {
	t.Parallel()
	finalPath := filepath.Join(t.TempDir(), "only.txt")
	RemoveTmpSafeFile(finalPath)
}

func TestSafeWriteFile_stringContent(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	target := filepath.Join(dir, "msg.txt")
	if err := SafeWriteFile(target, "hello-string", fs.FileMode(0o644)); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello-string" {
		t.Fatalf("got %q", got)
	}
}

func TestSafeWriteFile_emptyString(t *testing.T) {
	t.Parallel()
	target := filepath.Join(t.TempDir(), "empty.txt")
	if err := SafeWriteFile(target, "", fs.FileMode(0o600)); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("want empty file, got %d bytes", len(got))
	}
}
