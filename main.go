package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/noblemajo/deployit/lib/dit"
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

	hosts := []dit.SshTaskHost{}

	i = 0
	for connecitonUrl, sshTaskList := range sshTasks {
		taskHost, err := dit.NewSshTaskHost(
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
