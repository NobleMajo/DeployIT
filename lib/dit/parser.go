package dit

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/noblemajo/deployit/lib/sshutils"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SshTaskHost struct {
	ID            int
	ConnecitonUrl string
	SshConfig     sshutils.SshConfig
	Tasks         []Task
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

	tasks := []Task{}
	var newTask Task

	for _, rawTask := range rawTasks {
		newTask, err = ParseTask(rawTask)
		if err != nil {
			return SshTaskHost{}, err
		}

		tasks = append(tasks, newTask)
	}

	return SshTaskHost{
		ID:            id,
		ConnecitonUrl: connecitonUrl,
		SshConfig:     sshConfig,
		Tasks:         tasks,
	}, nil
}

func (taskHost *SshTaskHost) PrecheckAll() error {
	for _, task := range taskHost.Tasks {
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
		taskHost.SshConfig,
		func(
			sftp *sftp.Client,
			session *ssh.Session,
		) error {
			for id, task := range taskHost.Tasks {
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
