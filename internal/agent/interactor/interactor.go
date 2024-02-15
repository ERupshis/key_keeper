package interactor

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/utils"
	"github.com/erupshis/key_keeper/internal/common/logger"
	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
)

type Interactor struct {
	rd   *Reader
	wr   *Writer
	logs logger.BaseLogger
}

func NewInteractor(rd *Reader, wr *Writer, logs logger.BaseLogger) *Interactor {
	return &Interactor{
		rd:   rd,
		wr:   wr,
		logs: logs,
	}
}

func (i *Interactor) Printf(format string, a ...any) int {
	defer deferutils.ExecWithLogError(i.wr.Flush, i.logs)

	n, err := fmt.Fprintf(i.wr, format, a...)
	if err != nil {
		i.logs.Infof("print data in writer: %v", err)
	}
	return n
}

func (i *Interactor) Writer() *Writer {
	return i.wr
}

func (i *Interactor) ReadCommand() ([]string, bool) {
	i.Printf("enter command (or '%s'): ", utils.CommandExit)
	command, _, err := i.GetUserInputAndValidate(nil)
	if errors.Is(err, io.EOF) {
		return []string{utils.CommandExit}, true
	}

	command = strings.TrimSpace(command)
	commandParts := strings.Split(command, " ")
	if len(commandParts) == 0 || command == "" {
		i.Printf("Empty command. Try again\n")
		return nil, false
	}

	return commandParts, true
}

func (i *Interactor) GetUserInputAndValidate(regex *regexp.Regexp) (string, bool, error) {
	input, err := i.rd.getUserInput()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return "", true, errors.Join(io.EOF, errs.ErrInterruptedByUser)
		}

		i.Printf("unexpected error, try again: %v\n", err)
		return "", false, nil
	}

	input = strings.TrimSpace(strings.TrimRight(input, "\r\n"))
	if input == utils.CommandCancel {
		return input, true, errs.ErrInterruptedByUser
	}

	if regex != nil && !regex.MatchString(input) {
		i.Printf("incorrect input, try again or interrupt by '%s' command: ", utils.CommandCancel)
		return "", false, nil
	}

	return input, true, nil
}
