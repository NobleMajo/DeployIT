package stringfs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// used to get the relative paths of the home directory and the current working directory.
func GetRelativePaths() (string, string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", "", errors.New("cant get current working dir:\n> " + err.Error())
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", errors.New("cant get users home dir:\n> " + err.Error())
	}

	return homeDir, wd, nil
}

// used when a string pointer needs to be parsed to an absolute path.
func AbsolutePathRef(path *string) error {
	homeDir, workDir, err := GetRelativePaths()
	if err != nil {
		return err
	}

	return AbsolutePathRefFrom(homeDir, workDir, path)
}

// used if parsing a path for a different relative base, work dir, user or host.
func AbsolutePathRefFrom(homeDir string, workDir string, path *string) error {
	if path == nil {
		return errors.New("path is nil")
	}

	*path = strings.TrimSpace(*path)

	if strings.HasPrefix(*path, "~") {
		*path = strings.Replace(*path, "~", homeDir, 1)
	}

	if !strings.HasPrefix(*path, "/") {
		*path = workDir + "/" + *path
	}

	*path = filepath.Join(*path)

	return nil
}

// used when a string needs to be parsed to an absolute path.
func AbsolutePath(path string) (string, error) {
	homeDir, workDir, err := GetRelativePaths()
	if err != nil {
		return "", err
	}

	return AbsolutePathFrom(homeDir, workDir, path)
}

// used if parsing a path for a different relative base, work dir, user or host.
func AbsolutePathFrom(homeDir string, workDir string, path string) (string, error) {
	path = strings.TrimSpace(path)

	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", homeDir, 1)
	}

	if !strings.HasPrefix(path, "/") {
		path = workDir + "/" + path
	}

	path = filepath.Join(path)

	return path, nil
}
