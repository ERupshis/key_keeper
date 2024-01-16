package interactor

import (
	"bufio"
	"fmt"
	"os"
)

type Reader struct {
	*bufio.Reader
}

func NewReader(stream *os.File) *Reader {
	return &Reader{
		Reader: bufio.NewReader(stream),
	}
}

func (r *Reader) getUserInput() (string, error) {
	input, err := r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}

	return input, nil
}
