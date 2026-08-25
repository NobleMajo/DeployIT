package dit

import (
	"strings"
	"testing"
)

func TestNewSshTaskHost_success(t *testing.T) {
	t.Parallel()
	h, err := NewSshTaskHost(0, "ssh://alice@example.com/app!secret", []string{"CMD@exit 0"})
	if err != nil {
		t.Fatal(err)
	}
	if h.ID != 0 || h.ConnectionURL != "ssh://alice@example.com/app!secret" {
		t.Fatalf("unexpected host: %+v", h)
	}
	if len(h.Tasks) != 1 {
		t.Fatalf("tasks: %d", len(h.Tasks))
	}
	if _, ok := h.Tasks[0].(*CommandTask); !ok {
		t.Fatalf("want CommandTask, got %T", h.Tasks[0])
	}
}

func TestNewSshTaskHost_sshConfigFromURL(t *testing.T) {
	t.Parallel()
	h, err := NewSshTaskHost(0, "ssh://bob@deploy.example:2222/var/www!mypass", []string{"CMD@true"})
	if err != nil {
		t.Fatal(err)
	}
	cfg := h.SshConfig
	if cfg.User != "bob" || cfg.Host != "deploy.example" || cfg.Port != 2222 {
		t.Fatalf("unexpected ssh config: %+v", cfg)
	}
	if cfg.Password != "mypass" {
		t.Fatalf("Password = %q", cfg.Password)
	}
	if cfg.TargetDir != "/var/www" {
		t.Fatalf("TargetDir = %q", cfg.TargetDir)
	}
}

func TestNewSshTaskHost_multipleTasks(t *testing.T) {
	t.Parallel()
	h, err := NewSshTaskHost(0, "ssh://u@h.example/x!p", []string{
		"CMD@date",
		"CMD@whoami",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(h.Tasks) != 2 {
		t.Fatalf("want 2 tasks, got %d", len(h.Tasks))
	}
	for i, task := range h.Tasks {
		if _, ok := task.(*CommandTask); !ok {
			t.Fatalf("task %d: want CommandTask, got %T", i, task)
		}
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
	if !strings.Contains(pcheckErr.Error(), "'0'") {
		t.Fatalf("expected host id in message: %v", pcheckErr)
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
