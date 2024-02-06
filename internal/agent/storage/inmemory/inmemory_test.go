package inmemory

import (
	"testing"
	"time"

	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	type args struct {
		cryptHasher *ska.SKA
	}
	type want struct {
		errOccur assert.ComparisonAssertionFunc
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "base",
			args: args{
				cryptHasher: nil,
			},
			want: want{
				errOccur: assert.NotEqual,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewStorage(tt.args.cryptHasher)
			tt.want.errOccur(t, storage, nil)
		})
	}
}

func TestStorage_resetNextFreeIdx(t *testing.T) {
	type fields struct {
		records     []models.Record
		cryptHasher *ska.SKA
	}
	type want struct {
		idx      int64
		errOccur assert.ComparisonAssertionFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "base",
			fields: fields{
				records: []models.Record{
					{-1, models.Data{}, false, time.Now()},
					{-2, models.Data{}, false, time.Now()},
					{-3, models.Data{}, false, time.Now()},
				},
				cryptHasher: nil,
			},
			want: want{
				idx:      -3,
				errOccur: assert.Equal,
			},
		},
		{
			name: "only synced records",
			fields: fields{
				records: []models.Record{
					{1, models.Data{}, false, time.Now()},
					{2, models.Data{}, false, time.Now()},
					{3, models.Data{}, false, time.Now()},
				},
				cryptHasher: nil,
			},
			want: want{
				idx:      0,
				errOccur: assert.Equal,
			},
		},
		{
			name: "no records",
			fields: fields{
				records:     []models.Record{},
				cryptHasher: nil,
			},
			want: want{
				idx:      0,
				errOccur: assert.Equal,
			},
		},
		{
			name: "mixed records",
			fields: fields{
				records: []models.Record{
					{1, models.Data{}, false, time.Now()},
					{-2, models.Data{}, false, time.Now()},
					{3, models.Data{}, false, time.Now()},
				},
				cryptHasher: nil,
			},
			want: want{
				idx:      -2,
				errOccur: assert.Equal,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Storage{
				records:     tt.fields.records,
				cryptHasher: tt.fields.cryptHasher,
			}

			s.resetNextFreeIdx()
			tt.want.errOccur(t, tt.want.idx, s.freeIdx)
		})
	}
}

func TestStorage_getNextFreeIdx(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		{
			name: "base",
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{}
			assert.Equalf(t, tt.want, s.getNextFreeIdx(), "getNextFreeIdx()")
		})
	}
}
