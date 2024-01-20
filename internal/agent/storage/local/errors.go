package local

import (
	"fmt"
)

var (
	ErrFileIsNotOpen       = fmt.Errorf("file is not open")
	ErrIncorrectPassPhrase = fmt.Errorf("incorrect passphrase")
)
