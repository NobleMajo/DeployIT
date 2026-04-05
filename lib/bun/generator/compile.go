package main

import (
	"os"
	"os/exec"
)

func bunCompile(wd string, bunPath string, entrypoint string, outfile string, pipeStreams bool) error {
	_ = os.Remove(outfile)

	cmd := exec.Command(bunPath, "build", "--compile", entrypoint, "--outfile", outfile)
	cmd.Dir = wd

	if pipeStreams {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}
