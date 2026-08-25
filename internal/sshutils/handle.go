package sshutils

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func HandleSftp(
	sshConfig SshConfig,
	handle func(
		*sftp.Client,
		*ssh.Session,
	) error,
) error {
	err := sshConfig.VerifySshConfig()
	if err != nil {
		return errors.New("error verifying ssh config: " + err.Error())
	}

	authMethods := []ssh.AuthMethod{}

	if len(sshConfig.Password) > 0 {
		authMethods = append(authMethods, ssh.Password(sshConfig.Password))
	}

	if len(sshConfig.PrivateKey) > 0 {
		signer, err := ssh.ParsePrivateKey([]byte(sshConfig.PrivateKey))
		if err != nil {
			return errors.New("error parsing private key: " + err.Error())
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.New("error getting home directory: " + err.Error())
	}

	hostKeyCallback, err := knownhosts.New(filepath.Join(homeDir, ".ssh", "known_hosts"))
	if err != nil {
		return errors.New("error parsing known hosts: " + err.Error())
	}

	conf := &ssh.ClientConfig{
		User:            sshConfig.User,
		HostKeyCallback: hostKeyCallback,
		Auth:            authMethods,
	}

	// sftp
	sftpSshClient, err := ssh.Dial("tcp", sshConfig.Host+":"+strconv.Itoa(sshConfig.Port), conf)
	if err != nil {
		return errors.New("error dialing: " + err.Error())
	}
	defer sftpSshClient.Close()

	sftp, err := sftp.NewClient(
		sftpSshClient,
	)
	if err != nil {
		return errors.New("error creating sftp client: " + err.Error())
	}
	defer sftp.Close()

	// session
	sessionSshClient, err := ssh.Dial("tcp", sshConfig.Host+":"+strconv.Itoa(sshConfig.Port), conf)
	if err != nil {
		return errors.New("error dialing: " + err.Error())
	}
	defer sessionSshClient.Close()

	session, err := sessionSshClient.NewSession()
	if err != nil {
		return errors.New("error creating ssh session: " + err.Error())
	}
	defer session.Close()

	// handle
	err = handle(sftp, session)
	if err != nil {
		return errors.New("error handling: " + err.Error())
	}

	return nil
}
