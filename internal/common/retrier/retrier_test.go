package retrier

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
)

// DatabaseErrorsToRetry errors to retry request to database.
var DatabaseErrorsToRetry = []error{
	errors.New(pgerrcode.UniqueViolation),
	errors.New(pgerrcode.ConnectionException),
	errors.New(pgerrcode.ConnectionDoesNotExist),
	errors.New(pgerrcode.ConnectionFailure),
	errors.New(pgerrcode.SQLClientUnableToEstablishSQLConnection),
	errors.New(pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection),
	errors.New(pgerrcode.TransactionResolutionUnknown),
	errors.New(pgerrcode.ProtocolViolation),
}

func Test_canRetryCall(t *testing.T) {
	type args struct {
		err              error
		repeatableErrors []error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid",
			args: args{
				err:              errors.New(`08000`),
				repeatableErrors: DatabaseErrorsToRetry,
			},
			want: true,
		},
		{
			name: "valid with missing slice",
			args: args{
				err:              errors.New(`any error`),
				repeatableErrors: nil,
			},
			want: true,
		},
		{
			name: "invalid error is not in slice",
			args: args{
				err:              errors.New(`any error`),
				repeatableErrors: DatabaseErrorsToRetry,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, canRetryCall(tt.args.err, tt.args.repeatableErrors))
		})
	}
}

func TestRetryCallWithTimeout(t *testing.T) {
	type args struct {
		ctx              context.Context
		intervals        []int
		repeatableErrors []error
		callback         func(context.Context) (interface{}, error)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				ctx:              context.Background(),
				intervals:        []int{1, 1, 1},
				repeatableErrors: nil,
				callback: func(ctx context.Context) (interface{}, error) {
					<-ctx.Done()
					return nil, errors.New(pgerrcode.ConnectionException)

				},
			},
			wantErr: true,
		},
		{
			name: "valid with success",
			args: args{
				ctx:              context.Background(),
				intervals:        nil,
				repeatableErrors: nil,
				callback: func(ctx context.Context) (interface{}, error) {
					return nil, nil
				},
			},
			wantErr: false,
		},
		{
			name: "valid should retry",
			args: args{
				ctx:              context.Background(),
				intervals:        []int{1, 1, 1},
				repeatableErrors: DatabaseErrorsToRetry,
				callback: func(ctx context.Context) (interface{}, error) {
					<-ctx.Done()
					return nil, errors.New(pgerrcode.ConnectionException)
				},
			},
			wantErr: true,
		},
		{
			name: "valid shouldn't retry",
			args: args{
				ctx:              context.Background(),
				intervals:        []int{1, 1, 1},
				repeatableErrors: DatabaseErrorsToRetry,
				callback: func(ctx context.Context) (interface{}, error) {
					<-ctx.Done()
					return nil, errors.New("some error")
				},
			},
			wantErr: true,
		},
	}
	for _, ttB := range tests {
		tt := ttB
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := RetryCallWithTimeout(tt.args.ctx, tt.args.intervals, tt.args.repeatableErrors, tt.args.callback)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
