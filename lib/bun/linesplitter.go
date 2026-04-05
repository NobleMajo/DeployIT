package bun

import (
	"bytes"
	"fmt"
)

const maxLineSize = 16 << 20 // 16MB

type LineSplitter struct {
	onLine  func([]byte)
	buf     []byte
	maxSize int
}

func NewLineSplitter(onLine func([]byte)) *LineSplitter {
	return &LineSplitter{
		onLine:  onLine,
		maxSize: maxLineSize,
	}
}

func (w *LineSplitter) Write(p []byte) (int, error) {
	if len(w.buf)+len(p) > maxLineSize {
		return 0, fmt.Errorf("worker stream line exceeds max size (%d bytes)", maxLineSize)
	}
	w.buf = append(w.buf, p...)
	for {
		idx := bytes.IndexByte(w.buf, '\n')
		if idx < 0 {
			break
		}
		line := append([]byte(nil), w.buf[:idx]...)
		w.onLine(line)
		w.buf = append(w.buf[:0], w.buf[idx+1:]...)
	}
	return len(p), nil
}

func (w *LineSplitter) flush() {
	if len(bytes.TrimSpace(w.buf)) == 0 {
		w.buf = w.buf[:0]
		return
	}
	line := append([]byte(nil), w.buf...)
	w.onLine(line)
	w.buf = w.buf[:0]
}
