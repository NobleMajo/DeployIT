package stringfs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func ParsePathRef(path *string) error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.New("cant get current working dir:\n> " + err.Error())
	}

	return ParsePathRefFrom(path, wd)
}

func ParsePathRefFrom(path *string, cwd string) error {
	if path == nil {
		return errors.New("path is nil")
	}

	*path = strings.TrimSpace(*path)

	if strings.HasPrefix(*path, "~") {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return errors.New("cant get users home dir:\n> " + err.Error())
		}

		*path = strings.Replace(*path, "~", userHome, 1)
	}

	if !strings.HasPrefix(*path, "/") {
		*path = cwd + "/" + *path
	}

	*path = filepath.Join(*path)

	return nil
}

func ParsePath(path string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.New("cant get current working dir:\n> " + err.Error())
	}

	return ParsePathFrom(path, wd)
}

func ParsePathFrom(path string, cwd string) (string, error) {
	path = strings.TrimSpace(path)

	if strings.HasPrefix(path, "~") {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New("cant get users home dir:\n> " + err.Error())
		}

		path = strings.Replace(path, "~", userHome, 1)
	}

	if !strings.HasPrefix(path, "/") {
		path = cwd + "/" + path
	}

	path = filepath.Join(path)

	return path, nil
}
