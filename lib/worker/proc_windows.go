//go:build windows

package worker

import "os/exec"

func ApplySubprocessSysProcAttr(cmd *exec.Cmd) {}

func TerminateChild(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	_ = cmd.Process.Kill()
}
