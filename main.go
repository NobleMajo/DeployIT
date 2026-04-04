package main

import (
	"errors"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	dit "github.com/noblemajo/deployit/internal"
	"github.com/noblemajo/deployit/lib/sshutils"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var DisplayName string = "Unset"
var ShortName string = "unset"
var Version string = "?.?.?"
var Commit string = "???????"

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	slog.Info("starting",
		"display_name", DisplayName,
		"version", Version,
		"commit", Commit,
	)

	err := godotenv.Load()
	if err == nil {
		slog.Info("environment variables loaded", "source", ".env")
	}

	sshTasks := map[string][]string{}

	var j int
	var i int = 0
	for {
		connecitonUrl := os.Getenv("DIT_NODE" + strconv.Itoa(i+1))

		if connecitonUrl == "" {
			if i == 0 {
				slog.Error("no ssh config for node", "node", i+1)
				os.Exit(1)
			}

			break
		}

		sshTasks[connecitonUrl] = []string{}

		j = 0
		for {
			sshTask := os.Getenv("DIT_NODE" + strconv.Itoa(i+1) + "_TASK" + strconv.Itoa(j+1))

			if sshTask == "" {
				if j == 0 {
					slog.Error("no ssh task for node", "node", i+1, "task_index", j+1)
					os.Exit(1)
				}

				break
			}

			sshTasks[connecitonUrl] = append(sshTasks[connecitonUrl], sshTask)

			j++
		}

		i++
	}

	hosts := []SshTaskHost{}

	i = 0
	for connecitonUrl, sshTaskList := range sshTasks {
		taskHost, err := NewSshTaskHost(
			i,
			connecitonUrl,
			sshTaskList,
		)

		if err != nil {
			slog.Error("ssh task host setup failed", "err", err)
			os.Exit(1)
		}

		hosts = append(hosts, taskHost)

		i++
	}

	for _, host := range hosts {
		err := host.PrecheckAll()
		if err != nil {
			slog.Error("precheck failed", "err", err)
			os.Exit(1)
		}
	}

	for _, host := range hosts {
		err := host.Deploy()
		if err != nil {
			slog.Error("deploy failed", "err", err)
			os.Exit(1)
		}
	}

	slog.Info("done")
}

type SshTaskHost struct {
	ID            int
	connecitonUrl string
	sshConfig     sshutils.SshConfig
	tasks         []dit.Task
}

func NewSshTaskHost(
	id int,
	connecitonUrl string,
	rawTasks []string,
) (SshTaskHost, error) {
	sshConfig, err := sshutils.NewSshConfig(connecitonUrl)
	if err != nil {
		return SshTaskHost{}, err
	}

	tasks := []dit.Task{}
	var newTask dit.Task

	for _, rawTask := range rawTasks {
		newTask, err = dit.ParseTask(rawTask)
		if err != nil {
			return SshTaskHost{}, err
		}

		tasks = append(tasks, newTask)
	}

	return SshTaskHost{
		ID:            id,
		connecitonUrl: connecitonUrl,
		sshConfig:     sshConfig,
		tasks:         tasks,
	}, nil
}

func (taskHost *SshTaskHost) PrecheckAll() error {
	for _, task := range taskHost.tasks {
		err := task.Precheck()
		if err != nil {
			return errors.New(
				"precheck failed for '" + strconv.Itoa(taskHost.ID) +
					"' task '" + task.Raw() + "': " +
					err.Error(),
			)
		}
	}

	return nil
}

func (taskHost *SshTaskHost) Deploy() error {
	return sshutils.HandleSftp(
		taskHost.sshConfig,
		func(
			sftp *sftp.Client,
			session *ssh.Session,
		) error {
			for id, task := range taskHost.tasks {
				slog.Info("execute task", "task", task.Raw())
				err := task.Execute(sftp, session)
				if err != nil {
					return errors.New(
						"error host-" + strconv.Itoa(taskHost.ID) +
							" executing task-" + strconv.Itoa(id) + ": " +
							err.Error(),
					)
				}
			}

			return nil
		},
	)
}
