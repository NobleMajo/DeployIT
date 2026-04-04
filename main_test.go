package main

import (
	"strings"
	"testing"

	dit "github.com/noblemajo/deployit/internal"
)

func TestNewSshTaskHost_success(t *testing.T) {
	t.Parallel()
	h, err := NewSshTaskHost(0, "ssh://alice@example.com/app!secret", []string{"CMD@exit 0"})
	if err != nil {
		t.Fatal(err)
	}
	if h.ID != 0 || h.connecitonUrl != "ssh://alice@example.com/app!secret" {
		t.Fatalf("unexpected host: %+v", h)
	}
	if len(h.tasks) != 1 {
		t.Fatalf("tasks: %d", len(h.tasks))
	}
	if _, ok := h.tasks[0].(*dit.CommandTask); !ok {
		t.Fatalf("want CommandTask, got %T", h.tasks[0])
	}
}

func TestNewSshTaskHost_invalidTask(t *testing.T) {
	t.Parallel()
	_, err := NewSshTaskHost(0, "ssh://bob@host.example/x!pw", []string{"UPLOAD@only"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid upload task") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestNewSshTaskHost_invalidConnectionString(t *testing.T) {
	t.Parallel()
	_, err := NewSshTaskHost(0, "ssh://alice@server.example/app/path", []string{"CMD@true"})
	if err == nil {
		t.Fatal("expected error: connection string must include ! or * credentials")
	}
	if !strings.Contains(err.Error(), "invalid path credentials") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestPrecheckAll_emptyTasks(t *testing.T) {
	t.Parallel()
	h, err := NewSshTaskHost(1, "ssh://u@h.example/z!p", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := h.PrecheckAll(); err != nil {
		t.Fatal(err)
	}
}

func TestPrecheckAll_uploadMissingFile(t *testing.T) {
	t.Parallel()
	h, err := NewSshTaskHost(0, "ssh://u@h.example/dir!pw", []string{"UPLOAD@/this/path/does/not/exist@/remote"})
	if err != nil {
		t.Fatal(err)
	}
	pcheckErr := h.PrecheckAll()
	if pcheckErr == nil {
		t.Fatal("expected precheck error")
	}
	if !strings.Contains(pcheckErr.Error(), "precheck failed") {
		t.Fatalf("unexpected: %v", pcheckErr)
	}
}

func TestPrecheckAll_secondTaskFails(t *testing.T) {
	t.Parallel()
	h, err := NewSshTaskHost(0, "ssh://u@h.example/z!pw", []string{
		"CMD@true",
		"UPLOAD@/no/such/local/file@/remote",
	})
	if err != nil {
		t.Fatal(err)
	}
	pcheckErr := h.PrecheckAll()
	if pcheckErr == nil {
		t.Fatal("expected precheck error on second task")
	}
	if !strings.Contains(pcheckErr.Error(), "precheck failed") {
		t.Fatalf("unexpected: %v", pcheckErr)
	}
}
