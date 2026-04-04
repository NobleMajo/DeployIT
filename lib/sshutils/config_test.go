package sshutils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewSshConfig_passwordOnly(t *testing.T) {
	t.Parallel()
	raw := "ssh://alice@example.com:2222/var/app!s3cret"
	cfg, err := NewSshConfig(raw)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.User != "alice" || cfg.Host != "example.com" || cfg.Port != 2222 {
		t.Fatalf("unexpected identity: %+v", cfg)
	}
	if cfg.TargetDir != "/var/app" {
		t.Fatalf("TargetDir = %q", cfg.TargetDir)
	}
	if cfg.Password != "s3cret" || cfg.PrivateKey != "" {
		t.Fatalf("password/key mismatch: password=%q keyEmpty=%v", cfg.Password, cfg.PrivateKey == "")
	}
}

func TestNewSshConfig_privateKeyFromFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "id_rsa")
	content := "fake-pem-key-material\n"
	if err := os.WriteFile(keyPath, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	raw := "ssh://bob@host.example/x/y*" + keyPath
	cfg, err := NewSshConfig(raw)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.User != "bob" || cfg.Host != "host.example" {
		t.Fatalf("unexpected: %+v", cfg)
	}
	if cfg.Password != "" {
		t.Fatalf("expected empty password, got %q", cfg.Password)
	}
	if cfg.PrivateKey != content {
		t.Fatalf("PrivateKey mismatch: got %q", cfg.PrivateKey)
	}
}

func TestNewSshConfig_missingCredentialDelimiter(t *testing.T) {
	t.Parallel()
	_, err := NewSshConfig("ssh://u@h:22/path")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid path credentials") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSshConfig_starAfterBang(t *testing.T) {
	t.Parallel()
	_, err := NewSshConfig("ssh://u@h/x!pass*key")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "privateKey needs to be defined before") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSshConfig_wrongScheme(t *testing.T) {
	t.Parallel()
	_, err := NewSshConfig("http://u@h/x!p")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid url scheme") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSshConfig_portOutOfRange(t *testing.T) {
	t.Parallel()
	_, err := NewSshConfig("ssh://u@h:70000/x!p")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid port") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSshConfig_starOnlyInlineKey(t *testing.T) {
	t.Parallel()
	cfg, err := NewSshConfig("ssh://carol@host.example/sub/path*my-inline-key-material")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Password != "" {
		t.Fatalf("expected empty password, got %q", cfg.Password)
	}
	if cfg.PrivateKey != "my-inline-key-material" {
		t.Fatalf("PrivateKey = %q", cfg.PrivateKey)
	}
	if cfg.Port != 22 {
		t.Fatalf("default port want 22, got %d", cfg.Port)
	}
}

func TestNewSshConfig_defaultPort22(t *testing.T) {
	t.Parallel()
	cfg, err := NewSshConfig("ssh://dave@server.example/var/www!pw")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != 22 {
		t.Fatalf("Port = %d", cfg.Port)
	}
}

func TestNewSshConfig_portNotNumeric(t *testing.T) {
	t.Parallel()
	_, err := NewSshConfig("ssh://u@example.com:abc/var!pw")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "port") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestNewSshConfig_emptyPrivateKeyFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "empty_key")
	if err := os.WriteFile(keyPath, []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := NewSshConfig("ssh://u@h.example/x*" + keyPath)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestNewSshConfig_starBangPasswordPreservesInnerBang(t *testing.T) {
	t.Parallel()
	raw := "ssh://u@example.com/var*mykeymaterial!part2!part3"
	cfg, err := NewSshConfig(raw)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PrivateKey != "mykeymaterial" {
		t.Fatalf("PrivateKey = %q", cfg.PrivateKey)
	}
	if cfg.Password != "part2!part3" {
		t.Fatalf("Password = %q", cfg.Password)
	}
}

func TestNewSshConfig_privateKeyFileMissing(t *testing.T) {
	t.Parallel()
	missing := filepath.Join(t.TempDir(), "missing-private-key-file")
	_, err := NewSshConfig("ssh://u@h.example/x*" + missing)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "reading private key") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestNewSshConfig_emptyUsernameRejected(t *testing.T) {
	t.Parallel()
	_, err := NewSshConfig("ssh://@example.com/var/www!secret")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "User is required") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestNewSshConfig_emptyHostnameRejected(t *testing.T) {
	t.Parallel()
	_, err := NewSshConfig("ssh://alice@/var/www/app!secret")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "host is required") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestVerifySshConfig_valid(t *testing.T) {
	t.Parallel()
	cfg := SshConfig{
		User: "u", Host: "h", Port: 22,
		Password: "x", PrivateKey: "",
	}
	if err := cfg.VerifySshConfig(); err != nil {
		t.Fatal(err)
	}
}

func TestVerifySshConfig_missingAuth(t *testing.T) {
	t.Parallel()
	cfg := SshConfig{User: "u", Host: "h", Port: 22}
	if err := cfg.VerifySshConfig(); err == nil {
		t.Fatal("expected error")
	}
}

func TestVerifySshConfig_missingUser(t *testing.T) {
	t.Parallel()
	cfg := SshConfig{Host: "h", Port: 22, Password: "p"}
	if err := cfg.VerifySshConfig(); err == nil {
		t.Fatal("expected error")
	}
}

func TestVerifySshConfig_missingHost(t *testing.T) {
	t.Parallel()
	cfg := SshConfig{User: "u", Port: 22, Password: "p"}
	if err := cfg.VerifySshConfig(); err == nil {
		t.Fatal("expected error")
	}
}

func TestVerifySshConfig_portZero(t *testing.T) {
	t.Parallel()
	cfg := SshConfig{User: "u", Host: "h", Port: 0, Password: "p"}
	vErr := cfg.VerifySshConfig()
	if vErr == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(vErr.Error(), "invalid port") {
		t.Fatalf("unexpected: %v", vErr)
	}
}

func TestVerifySshConfig_portTooHigh(t *testing.T) {
	t.Parallel()
	cfg := SshConfig{User: "u", Host: "h", Port: 99999, Password: "p"}
	vErr := cfg.VerifySshConfig()
	if vErr == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(vErr.Error(), "invalid port") {
		t.Fatalf("unexpected: %v", vErr)
	}
}
