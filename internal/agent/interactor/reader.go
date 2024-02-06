package interactor

import (
	"bufio"
	"fmt"
	"io"
)

type Reader struct {
	*bufio.Reader
}

func NewReader(reader io.Reader) *Reader {
	return &Reader{
		Reader: bufio.NewReader(reader),
	}
}

func (r *Reader) getUserInput() (string, error) {
	input, err := r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}

	return input, nil
}
