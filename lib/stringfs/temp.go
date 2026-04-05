package stringfs

import (
	"os"
	"runtime"
)

func CreateTemp(data []byte, name string) (path string, clean func(), err error) {
	pattern := name + "-*"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	tempFile, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", nil, err
	}

	path = tempFile.Name()
	_, err = tempFile.Write(data)
	if err != nil {
		_ = tempFile.Close()
		_ = os.Remove(path)
		return "", nil, err
	}

	err = tempFile.Close()
	if err != nil {
		_ = os.Remove(path)
		return "", nil, err
	}

	if runtime.GOOS != "windows" {
		err = os.Chmod(path, 0o755)
		if err != nil {
			_ = os.Remove(path)
			return "", nil, err
		}
	}

	return path, func() {
		_ = os.Remove(path)
	}, nil
}
