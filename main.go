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

	fastFailString := os.Getenv("DIT_FAST_FAIL")
	fastFail := fastFailString == "true"

	type sshNode struct {
		connectionURL string
		tasks         []string
	}
	var nodes []sshNode

	var j int
	i := 0
	for {
		connectionURL := os.Getenv("DIT_NODE" + strconv.Itoa(i+1))

		if connectionURL == "" {
			if i == 0 {
				if fastFail {
					slog.Error("no ssh config for node", "node", i+1)
					os.Exit(1)
				} else {
					slog.Info("no ssh config for node", "node", i+1)
					os.Exit(0)
				}
			}

			break
		}

		tasks := []string{}

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

			tasks = append(tasks, sshTask)

			j++
		}

		nodes = append(nodes, sshNode{connectionURL: connectionURL, tasks: tasks})

		i++
	}

	hosts := make([]dit.SshTaskHost, 0, len(nodes))
	for i, node := range nodes {
		taskHost, err := dit.NewSshTaskHost(
			i,
			node.connectionURL,
			node.tasks,
		)

		if err != nil {
			slog.Error("ssh task host setup failed", "err", err)
			os.Exit(1)
		}

		hosts = append(hosts, taskHost)
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
