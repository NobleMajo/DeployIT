package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	dit "coreunit.net/wgg/internal"
	"coreunit.net/wgg/lib/sshutils"
	"github.com/joho/godotenv"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var DisplayName string = "Unset"
var ShortName string = "unset"
var Version string = "?.?.?"
var Commit string = "???????"

func main() {
	fmt.Println(DisplayName + " version v" + Version + ", build " + Commit)

	err := godotenv.Load()
	if err == nil {
		fmt.Println("Environment variables from .env loaded")
	}

	sshTasks := map[string][]string{}

	var j int
	var i int = 0
	for {
		connecitonUrl := os.Getenv("DIT_NODE" + strconv.Itoa(i+1))

		if len(connecitonUrl) <= 0 {
			if i == 0 {
				log.Fatalln("no ssh config for node " + strconv.Itoa(i+1))
			}

			break
		}

		sshTasks[connecitonUrl] = []string{}

		j = 0
		for {
			sshTask := os.Getenv("DIT_NODE" + strconv.Itoa(i+1) + "_TASK" + strconv.Itoa(j+1))

			if len(sshTask) <= 0 {
				if j == 0 {
					log.Fatalln("no ssh task for node " + strconv.Itoa(i+1) + " and task " + strconv.Itoa(j+1))
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
			log.Fatalln(err)
		}

		hosts = append(hosts, taskHost)

		i++
	}

	for _, host := range hosts {
		err := host.PrecheckAll()
		if err != nil {
			log.Fatalln(err)
		}
	}

	for _, host := range hosts {
		err := host.Deploy()
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Println("done")
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
				fmt.Println("Execute task " + task.Raw())
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
