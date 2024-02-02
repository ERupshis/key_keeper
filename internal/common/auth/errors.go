package auth

import (
	"fmt"
)

var (
	ErrMismatchPassword = fmt.Errorf("password mismatch")
	ErrLoginOccupied    = fmt.Errorf("login already occupied")
)
