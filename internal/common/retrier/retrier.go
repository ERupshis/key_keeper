// Package retrier implements repeat requests logic.
package retrier

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// defIntervals default intervals for repeats.
var defIntervals = []int{1, 3, 5}

// RetryCallWithTimeout generates repeats of function call if error occurs.
// Args:
//   - ctx(context.Context),
//   - log Logs.BaseLogs,
//   - intervals([]int) - count of repeats and pause between them (secs.);
//   - repeatableErrors([]error) - errors - reasons to make repeat call. If empty - any error is signal to repeat call;
//   - callback (func(context.Context) error) - function to call.
func RetryCallWithTimeout[T *sql.Rows | interface{} | []byte](ctx context.Context, intervals []int, repeatableErrors []error,
	callback func(context.Context) (T, error)) (T, error) {
	var err error

	if intervals == nil {
		intervals = defIntervals
	}

	attemptNum := 0
	var errs []error

	for _, interval := range intervals {
		ctxWithTime, cancel := context.WithTimeout(ctx, time.Duration(interval)*time.Second)
		go waitContextToCancel(ctxWithTime, cancel, interval)

		var tmp T
		tmp, err = callback(ctxWithTime)
		if err == nil {
			return tmp, nil
		}

		attemptNum++
		errs = append(errs, fmt.Errorf("attempt #%d to perform action failed with error: %w", attemptNum, err))

		if !canRetryCall(err, repeatableErrors) {
			break
		}
	}

	var tmpFake T
	return tmpFake, errors.Join(errs...)
}

// waitContextToCancel goroutine to prevent timeout context leaking.
func waitContextToCancel(ctx context.Context, cancelFunc context.CancelFunc, interval int) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(time.Duration(interval) * time.Second):
		cancelFunc()
	}
}

// canRetryCall checks if generated error is in list of repeatableErrors.
func canRetryCall(err error, repeatableErrors []error) bool {
	if repeatableErrors == nil {
		return true
	}

	canRetry := false
	for _, repeatableError := range repeatableErrors {
		if err.Error() == repeatableError.Error() {
			canRetry = true
		}
	}

	return canRetry
}
