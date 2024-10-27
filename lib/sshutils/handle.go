package sshutils

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func HandleSftp(
	sshConfig SshConfig,
	handle func(
		*sftp.Client,
		*ssh.Session,
	) error,
) error {
	// var hostkeyCallback ssh.HostKeyCallback
	// hostkeyCallback, err = knownhosts.New(homeDir + "/.ssh/known_hosts")
	// if err != nil {
	// 	return errors.New("error parsing known hosts: " + err.Error())
	// }

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

	conf := &ssh.ClientConfig{
		User: sshConfig.User,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Auth: authMethods,
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

func JoinPath(sftp *sftp.Client, path ...string) (string, error) {
	if len(path) != 0 &&
		len(path[0]) != 0 {
		if strings.HasPrefix(path[0], "~/") ||
			strings.HasPrefix(path[0], "./") {
			cwd, err := sftp.Getwd()
			if err != nil {
				return "", errors.New("error getting cwd: " + err.Error())
			}

			if strings.HasPrefix(path[0], "../") {
				path[0] = cwd + "/" + path[0]
			} else {
				path[0] = cwd + path[0][1:]
			}
		}
	}

	return sftp.Join(path...), nil
}
