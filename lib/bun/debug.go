package bun

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
)

const maxLogLineSize = 512

func truncateWorkerLogLine(b []byte) string {
	if len(b) <= maxLogLineSize {
		return string(b)
	}
	return string(b[:maxLogLineSize]) + "…"
}

func parseJSONMessage(line []byte) (map[string]any, error) {
	var packet map[string]any
	err := json.Unmarshal(bytes.TrimSpace(line), &packet)
	return packet, err
}

func parseWorkerStdoutStream(ctx context.Context, log *slog.Logger, line []byte) {
	if len(bytes.TrimSpace(line)) == 0 {
		return
	}
	packet, err := parseJSONMessage(line)
	if err != nil {
		log.Log(ctx, slog.LevelError, "worker stdout: JSON parse error",
			"err", err, "raw", truncateWorkerLogLine(line))
		return
	}
	msg := "worker stdout"
	if s, ok := packet["data"].(string); ok && s != "" {
		msg = s
	}
	log.Log(ctx, slog.LevelInfo, msg, "stream", "stdout", "packet", packet)
}

func parseWorkerStderrStream(ctx context.Context, log *slog.Logger, line []byte) {
	if len(bytes.TrimSpace(line)) == 0 {
		return
	}
	packet, err := parseJSONMessage(line)
	if err != nil {
		log.Log(ctx, slog.LevelError, "worker stderr: JSON parse error",
			"err", err, "raw", truncateWorkerLogLine(line))
		return
	}
	msg := "worker stderr"
	if s, ok := packet["data"].(string); ok && s != "" {
		msg = s
	}
	log.Log(ctx, slog.LevelWarn, msg, "stream", "stderr", "packet", packet)
}
