// Package configutils provides functions for config handling.
package configutils

import (
	"fmt"
	"strconv"
	"time"
)

// SetEnvToParamIfNeed assigns environment value to param depends on param's type definition.
// Accepts *int64, *string as params, *time.Duration.
func SetEnvToParamIfNeed(param interface{}, val string) error {
	errMsg := "convert env variable '%s' to '%T'"
	if val == "" {
		if fmt.Sprint(param) == "" {
			return ErrMissingEnvValue
		}

		return nil
	}

	switch param := param.(type) {
	case *int64:
		{
			envVal, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return fmt.Errorf(errMsg+": %w", val, param, err)
			}

			*param = envVal
		}
	case *bool:
		{
			envVal, err := strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf(errMsg+": %w", val, param, err)
			}

			*param = envVal
		}
	case *time.Duration:
		{
			envVal, err := time.ParseDuration(val)
			if err != nil {
				return fmt.Errorf(errMsg+": %w", val, param, err)
			}

			*param = envVal
		}
	case *string:
		*param = val
	default:
		return ErrUnknownEnvType
	}

	return nil
}
