package interactor

import (
	"bufio"
	"os"
)

type Writer struct {
	*bufio.Writer
}

func NewWriter(stream *os.File) *Writer {
	return &Writer{
		Writer: bufio.NewWriter(stream),
	}
}
