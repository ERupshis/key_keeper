package errs

import (
	"fmt"
)

const (
	ErrProcessMsgBody = "process '%s' command: %w"
)

var (
	ErrIncorrectRecordType = fmt.Errorf("incorrect record type")
	ErrInterruptedByUser   = fmt.Errorf("interrupted by user")
	ErrUnexpected          = fmt.Errorf("unexpected error")
	ErrFilePathIsNotAbs    = fmt.Errorf("entered not absolute path")
)
