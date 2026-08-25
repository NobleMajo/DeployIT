package dit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseTask_upload(t *testing.T) {
	t.Parallel()
	task, err := ParseTask("UPLOAD@/local/file@/remote/file")
	if err != nil {
		t.Fatal(err)
	}
	u, ok := task.(*UploadTask)
	if !ok {
		t.Fatalf("want *UploadTask, got %T", task)
	}
	if u.FromPath != "/local/file" || u.ToPath != "/remote/file" {
		t.Fatalf("%+v", u)
	}
}

func TestParseTask_download(t *testing.T) {
	t.Parallel()
	task, err := ParseTask("DOWNLOAD@/r/src@/l/dst")
	if err != nil {
		t.Fatal(err)
	}
	d, ok := task.(*DownloadTask)
	if !ok {
		t.Fatalf("want *DownloadTask, got %T", task)
	}
	if d.FromPath != "/r/src" || d.ToPath != "/l/dst" {
		t.Fatalf("%+v", d)
	}
}

func TestParseTask_cmdWithSpaces(t *testing.T) {
	t.Parallel()
	task, err := ParseTask("CMD@echo hello world")
	if err != nil {
		t.Fatal(err)
	}
	c, ok := task.(*CommandTask)
	if !ok {
		t.Fatalf("want *CommandTask, got %T", task)
	}
	if c.Cmd != "echo hello world" {
		t.Fatalf("Cmd = %q", c.Cmd)
	}
}

func TestParseTask_empty(t *testing.T) {
	t.Parallel()
	_, err := ParseTask("")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseTask_uploadBadArity(t *testing.T) {
	t.Parallel()
	_, err := ParseTask("UPLOAD@only")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid upload task") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestParseTask_unknown(t *testing.T) {
	t.Parallel()
	_, err := ParseTask("OTHER@x")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown task") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestParseTask_downloadBadArity(t *testing.T) {
	t.Parallel()
	_, err := ParseTask("DOWNLOAD@only")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid download task") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestParseTask_cmdBadArity(t *testing.T) {
	t.Parallel()
	_, err := ParseTask("CMD@a@b")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid command task") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestParseTask_uploadTooManySegments(t *testing.T) {
	t.Parallel()
	_, err := ParseTask("UPLOAD@a@b@c")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid upload task") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestParseTask_downloadTooManySegments(t *testing.T) {
	t.Parallel()
	_, err := ParseTask("DOWNLOAD@a@b@c")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid download task") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestUploadTask_Precheck_emptyFrom(t *testing.T) {
	t.Parallel()
	task := &UploadTask{FromPath: ""}
	if err := task.Precheck(); err == nil {
		t.Fatal("expected error")
	}
}

func TestUploadTask_Precheck_regularFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	p := filepath.Join(dir, "data.bin")
	if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	task := &UploadTask{FromPath: p, ToPath: "/remote"}
	if err := task.Precheck(); err != nil {
		t.Fatal(err)
	}
}

func TestUploadTask_Precheck_notRegularFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	task := &UploadTask{FromPath: dir, ToPath: "/remote"}
	if err := task.Precheck(); err == nil {
		t.Fatal("expected error for directory path")
	}
}

func TestDownloadTask_Precheck_parentExists(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	to := filepath.Join(dir, "localfile")
	task := &DownloadTask{FromPath: "/remote/x", ToPath: to}
	if err := task.Precheck(); err != nil {
		t.Fatal(err)
	}
}

func TestDownloadTask_Precheck_parentMissing(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	to := filepath.Join(dir, "nope", "missing", "file")
	task := &DownloadTask{FromPath: "/r", ToPath: to}
	if err := task.Precheck(); err == nil {
		t.Fatal("expected error")
	}
}

func TestDownloadTask_Precheck_parentIsFileNotDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	fileAsParent := filepath.Join(dir, "not-a-dir")
	if err := os.WriteFile(fileAsParent, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	to := filepath.Join(fileAsParent, "out.txt")
	task := &DownloadTask{FromPath: "/remote/x", ToPath: to}
	pErr := task.Precheck()
	if pErr == nil {
		t.Fatal("expected error when parent exists but is not a directory")
	}
	if !strings.Contains(pErr.Error(), "local target dir is not a directory") {
		t.Fatalf("unexpected: %v", pErr)
	}
}
