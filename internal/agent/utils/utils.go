package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/erupshis/key_keeper/internal/agent/errs"
)

var reader = bufio.NewReader(os.Stdin)

func GetUserInputAndValidate(regex *regexp.Regexp) (string, bool, error) {
	input, err := getUserInput()
	input = strings.TrimRight(input, "\r\n")
	if err != nil {
		fmt.Printf("unexpected error, try again: %v\n", err)
		return "", false, nil
	}

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
