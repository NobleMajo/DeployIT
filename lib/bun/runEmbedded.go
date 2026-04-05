package bun

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/noblemajo/deployit/lib/stringfs"
	"github.com/noblemajo/deployit/lib/worker"
)

func RunEmbedded(ctx context.Context, log *slog.Logger) error {
	blob, err := EmbeddedBunBytes()
	if err != nil {
		return fmt.Errorf("read embedded bytes error: %w", err)
	}

	tempPath, clean, err := stringfs.CreateTemp(blob, "subprocess-blob")
	if err != nil {
		return fmt.Errorf("create temp file error: %w", err)
	}
	defer clean()

	stdoutW := NewLineSplitter(func(line []byte) {
		parseWorkerStdoutStream(ctx, log, line)
	})
	stderrW := NewLineSplitter(func(line []byte) {
		parseWorkerStderrStream(ctx, log, line)
	})
	defer func() {
		stdoutW.flush()
		stderrW.flush()
	}()

	w, err := worker.CreateWorker(worker.WorkerOptions{
		Ctx:        ctx,
		Cmd:        tempPath,
		PipeStdout: stdoutW,
		PipeStderr: stderrW,
		PipeStdin:  os.Stdin,
	})
	if err != nil {
		return fmt.Errorf("create worker error: %w", err)
	}

	err = <-w.Done
	if err != nil {
		return fmt.Errorf("worker done error: %w", err)
	}

	return nil
}
