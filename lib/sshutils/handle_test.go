package sshutils

import (
	"strings"
	"testing"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func TestHandleSftp_verifyFailsBeforeDial(t *testing.T) {
	t.Parallel()
	err := HandleSftp(SshConfig{}, func(*sftp.Client, *ssh.Session) error {
		t.Fatal("callback must not run when config is invalid")
		return nil
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "verifying ssh config") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestHandleSftp_invalidPrivateKeyBeforeDial(t *testing.T) {
	t.Parallel()
	cfg := SshConfig{
		User:       "u",
		Host:       "127.0.0.1",
		Port:       22,
		PrivateKey: "not-valid-openssh-pem",
	}
	err := HandleSftp(cfg, func(*sftp.Client, *ssh.Session) error {
		t.Fatal("callback must not run when key material is invalid")
		return nil
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "parsing private key") {
		t.Fatalf("unexpected: %v", err)
	}
}
