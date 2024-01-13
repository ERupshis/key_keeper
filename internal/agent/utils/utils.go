package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/erupshis/key_keeper/internal/agent/errs"
)

var reader = bufio.NewReader(os.Stdin)

// TODO: need Read/Writer class with DI streams
func GetUserInputAndValidate(regex *regexp.Regexp) (string, bool, error) {
	input, err := getUserInput()
	if err != nil {
		if !errors.Is(err, io.EOF) {
			fmt.Printf("unexpected error, try again: %v\n", err)
		}
		return "", false, nil
	}

	input = strings.TrimSpace(strings.TrimRight(input, "\r\n"))
	if input == CommandCancel {
		return input, true, errs.ErrInterruptedByUser
	}

	if !regex.MatchString(input) {
		fmt.Printf("incorrect input, try again or interrupt by 'cancel' command: ")
		return "", false, nil
	}

	return input, true, nil
}

func getUserInput() (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}

	return input, nil
}
