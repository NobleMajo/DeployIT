package dit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Task interface {
	Execute(
		sftp *sftp.Client,
		session *ssh.Session,
	) error
	Type() string
	Raw() string
	Precheck() error
}

type UploadTask struct {
	RawTask  string
	FromPath string
	ToPath   string
}

func (task *UploadTask) Type() string {
	return "UPLOAD"
}

func (task *UploadTask) Raw() string {
	return task.RawTask
}

func (task *UploadTask) Precheck() error {
	if task.FromPath == "" {
		return errors.New("upload source file is empty")
	}

	stats, err := os.Stat(task.FromPath)
	if err != nil {
		return err
	}

	if !stats.Mode().IsRegular() {
		return errors.New(task.FromPath + " is not a regular file")
	}

	return nil
}

func (task *UploadTask) Execute(
	sftp *sftp.Client,
	session *ssh.Session,
) error {
	srcFile, err := os.Open(task.FromPath)
	if err != nil {
		return errors.New("error opening source file: " + err.Error())
	}
	defer srcFile.Close()

	dstFile, err := sftp.Create(task.ToPath)
	if err != nil {
		return errors.New("error creating destination file: " + err.Error())
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errors.New("error copying file to remote: " + err.Error())
	}

	return nil
}

type DownloadTask struct {
	RawTask  string
	FromPath string
	ToPath   string
}

func (task *DownloadTask) Type() string {
	return "DOWNLOAD"
}

func (task *DownloadTask) Raw() string {
	return task.RawTask
}

func (task *DownloadTask) Precheck() error {
	if task.FromPath == "" {
		return errors.New("download source file is empty")
	}

	parentDir := filepath.Dir(task.ToPath)

	stats, err := os.Stat(parentDir)
	if err != nil {
		return errors.New("cant stat parent dir of local target: '" + parentDir + "': " + err.Error())
	}

	if !stats.IsDir() {
		return errors.New("local target dir is not a directory: '" + parentDir + "'")
	}

	return nil
}

func (task *DownloadTask) Execute(
	sftp *sftp.Client,
	session *ssh.Session,
) error {
	dstFile, err := os.Create(task.ToPath)
	if err != nil {
		return errors.New("error creating destination file: " + err.Error())
	}
	defer dstFile.Close()

	srcFile, err := sftp.Open(task.FromPath)
	if err != nil {
		return errors.New("error opening source file: " + err.Error())
	}
	defer srcFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errors.New("error copying file to local: " + err.Error())
	}

	return nil
}

type CommandTask struct {
	RawTask string
	Cmd     string
}

func (task *CommandTask) Type() string {
	return "CMD"
}

func (task *CommandTask) Raw() string {
	return task.RawTask
}

func (task *CommandTask) Precheck() error {
	return nil
}

func (task *CommandTask) Execute(
	sftp *sftp.Client,
	session *ssh.Session,
) error {
	out, err := session.CombinedOutput(task.Cmd)
	if err != nil {
		out2 := string(out)
		if len(out2) > 0 {
			return errors.New(
				"error executing command '" + task.Cmd +
					"': output: '" + string(out) +
					"', error: " + err.Error(),
			)
		} else {
			return errors.New(
				"error executing command '" + task.Cmd +
					"': empty output, error: " + err.Error(),
			)
		}
	}
	fmt.Println("\nCommand output of '" + task.Cmd + "':\n" + string(out) + "")

	return nil
}

func ParseTask(task string) (Task, error) {
	if task == "" {
		return nil, errors.New("cmd task command is empty")
	}

	splitted := strings.Split(task, "@")
	if splitted[0] == "UPLOAD" {
		if len(splitted) != 3 {
			return nil, errors.New(
				"invalid upload task: task has invalid format: " +
					"UPLOAD@<FromPath>@<ToPath> but is '" +
					task + "'",
			)
		}
		return &UploadTask{
			RawTask:  task,
			FromPath: splitted[1],
			ToPath:   splitted[2],
		}, nil
	} else if splitted[0] == "DOWNLOAD" {
		if len(splitted) != 3 {
			return nil, errors.New(
				"invalid download task: task has invalid format: " +
					"DOWNLOAD@<FromPath>@<ToPath> but is '" +
					task + "'",
			)
		}
		return &DownloadTask{
			RawTask:  task,
			FromPath: splitted[1],
			ToPath:   splitted[2],
		}, nil
	} else if splitted[0] == "CMD" {
		if len(splitted) != 2 {
			return nil, errors.New(
				"invalid command task: task has invalid format: " +
					"COMMAND@<Command> but is '" +
					task + "'",
			)
		}
		return &CommandTask{
			RawTask: task,
			Cmd:     splitted[1],
		}, nil
	} else {
		return nil, errors.New("cant parse task: unknown task: '" + task + "'")
	}
}
