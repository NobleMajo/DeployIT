package sshutils

import (
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"

	"coreunit.net/wgg/lib/stringfs"
)

type SshConfig struct {
	User       string
	Host       string
	Port       int
	TargetDir  string // used as relative dir for "./example" on the target server
	PrivateKey string // can be empty if Password is set
	Password   string // can be empty if PrivateKey is set
}

func NewSshConfig(
	rawConnecitonUrl string,
) (SshConfig, error) {
	var password string
	var privateKey string
	var connecitonUrl string

	// rawConnecitonUrl example "ssh://User@hostname:2222/path/to/directory@privateKeyPath!password"
	if strings.Contains(rawConnecitonUrl, "*") {
		if strings.Contains(rawConnecitonUrl, "!") {
			if strings.Index(rawConnecitonUrl, "*") > strings.Index(rawConnecitonUrl, "!") {
				return SshConfig{}, errors.New(
					"invalid path credentials, '*'-privateKey needs to be defined before '!'-password, as suffix '" +
						rawConnecitonUrl + "'",
				)
			}

			splitted := strings.Split(rawConnecitonUrl, "*")

			connecitonUrl = strings.TrimSpace(splitted[0])
			privateKey = strings.TrimSpace(strings.Join(splitted[1:], "*"))

			splitted = strings.Split(privateKey, "!")
			privateKey = strings.TrimSpace(splitted[0])
			password = strings.TrimSpace(strings.Join(splitted[1:], "!"))
		} else {
			splitted := strings.Split(rawConnecitonUrl, "*")
			connecitonUrl = strings.TrimSpace(splitted[0])
			password = ""
			privateKey = strings.TrimSpace(strings.Join(splitted[1:], "*"))
		}
	} else {
		if strings.Contains(rawConnecitonUrl, "!") {
			splitted := strings.Split(rawConnecitonUrl, "!")

			connecitonUrl = strings.TrimSpace(splitted[0])
			password = strings.TrimSpace(strings.Join(splitted[1:], "!"))
			privateKey = ""
		} else {
			return SshConfig{}, errors.New(
				"invalid path credentials, need '*' for privateKey " +
					"or '!' for password, is '" +
					rawConnecitonUrl + "'",
			)
		}
	}

	if len(privateKey) > 0 {
		if strings.HasPrefix(privateKey, "~/") ||
			strings.HasPrefix(privateKey, "./") ||
			strings.HasPrefix(privateKey, "../") ||
			strings.HasPrefix(privateKey, "/") ||
			strings.HasPrefix(privateKey, "file://") {
			privateKeyPath := &privateKey
			err := stringfs.ParsePathRef(privateKeyPath)
			if err != nil {
				return SshConfig{}, errors.New(
					"error parsing path: " +
						err.Error(),
				)
			}

			privateKeyBytes, err := os.ReadFile(*privateKeyPath)
			if err != nil {
				return SshConfig{}, errors.New(
					"error reading private key: " +
						err.Error(),
				)
			}

			privateKey = string(privateKeyBytes)

			if len(privateKey) <= 0 {
				return SshConfig{}, errors.New(
					"loaded private key from " +
						*privateKeyPath + " is empty",
				)
			}
		}
	}

	// connecitonUrl example "ssh://User@hostname:2222/path/to/directory"
	parsedURL, err := url.Parse(connecitonUrl)
	if err != nil {
		return SshConfig{}, errors.New("ssh url parse error: " + err.Error())
	}

	if parsedURL.Scheme != "ssh" {
		return SshConfig{}, errors.New("invalid url scheme, need 'ssh', is '" + parsedURL.Scheme + "'")
	}

	var port int
	portString := parsedURL.Port()
	if len(portString) == 0 {
		port = 22
	} else {
		port, err = strconv.Atoi(parsedURL.Port())
		if err != nil {
			return SshConfig{}, errors.New(
				"cant parse port to int, value '" +
					parsedURL.Port() + "': " +
					err.Error(),
			)
		}
	}

	if port < 1 || port > 65535 {
		return SshConfig{}, errors.New("invalid port, need 1-65535, is '" + parsedURL.Port() + "'")
	}

	sshConfig := SshConfig{
		User:       parsedURL.User.Username(),
		Host:       parsedURL.Hostname(),
		Port:       port,
		TargetDir:  parsedURL.Path,
		PrivateKey: privateKey,
		Password:   password,
	}

	err = sshConfig.VerifySshConfig()
	if err != nil {
		return SshConfig{}, err
	}

	return sshConfig, nil
}

func (sshConfig *SshConfig) VerifySshConfig() error {
	if sshConfig.Port < 1 || sshConfig.Port > 65535 {
		return errors.New("sshconfig verify error: invalid port number")
	}

	if len(sshConfig.Password) <= 0 && len(sshConfig.PrivateKey) <= 0 {
		return errors.New("sshconfig verify error: password or private key is required")
	}

	if len(sshConfig.User) <= 0 {
		return errors.New("sshconfig verify error: User is required")
	}

	if len(sshConfig.Host) <= 0 {
		return errors.New("sshconfig verify error: host is required")
	}

	return nil
}
