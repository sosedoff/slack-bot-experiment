package main

import (
	"strings"
)

type stdoutWriter struct {
	Lines chan string
}

func (w stdoutWriter) Write(buf []byte) (int, error) {
	line := strings.TrimSpace(string(buf))
	if line != "" {
		w.Lines <- line
	}
	return len(buf), nil
}
