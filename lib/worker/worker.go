package worker

import (
	"context"
	"errors"
	"io"
	"os/exec"
)

type WorkerOptions struct {
	Ctx        context.Context   // can be empty, fallback to background
	EnvVars    map[string]string // can be empty, fallback to empty/none
	EnvArray   []string          // can be empty, fallback to empty/none
	Cmd        string            // required not empty and min length 1
	Args       []string          // can be empty, fallback to empty/none
	PipeStdout io.Writer         // can be empty, fallback to dont pipe
	PipeStderr io.Writer         // can be empty, fallback to dont pipe
	PipeStdin  io.Reader         // can be empty, fallback to dont pipe
}

type Worker struct {
	Options WorkerOptions
	Cmd     *exec.Cmd
	Done    chan error
}

func CreateWorker(
	options WorkerOptions,
) (*Worker, error) {
	if options.Cmd == "" {
		return nil, errors.New("cmd is required and must be not empty")
	}
	if len(options.Args) < 1 {
		options.Args = []string{}
	}

	ctx := options.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	w := &Worker{
		Options: options,
		Done:    make(chan error, 1),
	}

	cmd := exec.CommandContext(ctx, options.Cmd, options.Args...)
	cmd.Env = parseEnvVars(options.EnvVars, options.EnvArray)
	ApplySubprocessSysProcAttr(cmd)

	if options.PipeStdout != nil {
		cmd.Stdout = options.PipeStdout
	}
	if options.PipeStderr != nil {
		cmd.Stderr = options.PipeStderr
	}
	if options.PipeStdin != nil {
		cmd.Stdin = options.PipeStdin
	}

	w.Cmd = cmd
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		w.Done <- cmd.Wait()
	}()

	return w, nil
}

func parseEnvVars(
	envVars map[string]string,
	envArray []string,
) []string {
	env := envArray
	for k, v := range envVars {
		env = append(env, k+"="+v)
	}
	return env
}
