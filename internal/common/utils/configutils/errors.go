package configutils

import (
	"fmt"
)

var (
	ErrMissingEnvValue = fmt.Errorf("missing value in env")
	ErrUnknownEnvType  = fmt.Errorf("unknown env param type")
)

func ErrCheckEnvsWrapper(err error) error {
	return fmt.Errorf("check environments: %w", err)
}
