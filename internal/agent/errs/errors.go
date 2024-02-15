package errs

import (
	"fmt"
)

const (
	ErrProcessMsgBody = "process '%s' command: %w"
)

var (
	ErrIncorrectRecordType       = fmt.Errorf("incorrect record type")
	ErrInterruptedByUser         = fmt.Errorf("interrupted by user")
	ErrUnexpected                = fmt.Errorf("unexpected error")
	ErrIncorrectServerActionType = fmt.Errorf("incorrect server action type")
)
