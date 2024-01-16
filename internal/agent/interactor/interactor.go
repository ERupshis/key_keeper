package interactor

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/utils"
)

type Interactor struct {
	rd *Reader
	wr *Writer
}

func NewInteractor(rd *Reader, wr *Writer) *Interactor {
	return &Interactor{
		rd: rd,
		wr: wr,
	}
}

func (i *Interactor) Printf(format string, a ...any) int {
	defer func() {
		_ = i.wr.Flush()
	}()

	n, _ := fmt.Fprintf(i.wr, format, a...) // TODO: handle error.

	return n
}

func (i *Interactor) Writer() *Writer {
	return i.wr
}

func (i *Interactor) ReadCommand() ([]string, bool) {
	i.Printf("Insert command (or '%s'): ", utils.CommandExit)
	command, _, _ := i.GetUserInputAndValidate(nil)
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
		if !errors.Is(err, io.EOF) {
			fmt.Printf("unexpected error, try again: %v\n", err)
		}
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
