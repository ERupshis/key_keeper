package local

import (
	"fmt"
)

var (
	ErrFileIsNotOpen = fmt.Errorf("file is not open")
	ErrDecryptData   = fmt.Errorf("decrypt data problem")
	ErrEncryptData   = fmt.Errorf("encrypt data problem")
)
