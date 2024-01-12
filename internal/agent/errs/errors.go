package errs

import (
	"fmt"
)

var (
	ErrIncorrectRecordType = fmt.Errorf("incorrect record type")
	ErrInterruptedByUser   = fmt.Errorf("interrupted by user")
)
