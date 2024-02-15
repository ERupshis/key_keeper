package interactor

import (
	"bufio"
	"io"
)

type Writer struct {
	*bufio.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		Writer: bufio.NewWriter(writer),
	}
}
