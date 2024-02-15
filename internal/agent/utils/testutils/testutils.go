package testutils

import (
	"bufio"
	"io"

	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/common/logger"
)

func CreateUserInteractor(rd io.Reader, wr io.Writer, logs logger.BaseLogger) *interactor.Interactor {
	reader := interactor.NewReader(bufio.NewReader(rd))
	writer := interactor.NewWriter(bufio.NewWriter(wr))
	return interactor.NewInteractor(reader, writer, logs)
}

func AddNewRow(in string) string {
	return in + "\n"
}
