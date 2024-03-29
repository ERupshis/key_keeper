package inmemory

import (
	"fmt"
	"reflect"
	"testing"

	models "github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/stretchr/testify/assert"
)

func TestStorage_RestoreRecords(t *testing.T) {
	type fields struct {
		records     []models.Record
		cryptHasher *ska.SKA
		freeIdx     int64
	}
	type args struct {
		records []models.Record
	}
	type want struct {
		err     assert.ErrorAssertionFunc
		idx     int64
		records []models.Record
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "base",
			fields: fields{
				records:     nil,
				cryptHasher: nil,
				freeIdx:     0,
			},
			args: args{
				records: []models.Record{
					{ID: 1},
					{ID: -2},
					{ID: 3},
				},
			},
			want: want{
				err: assert.NoError,
				idx: -2,
				records: []models.Record{
					{ID: 1},
					{ID: -2},
					{ID: 3},
				},
			},
		},
		{
			name: "already with records",
			fields: fields{
				records: []models.Record{
					{ID: -1},
					{ID: -2},
					{ID: -3},
				},
				cryptHasher: nil,
				freeIdx:     0,
			},
			args: args{
				records: []models.Record{
					{ID: 1},
					{ID: 2},
					{ID: 3},
				},
			},
			want: want{
				err: assert.NoError,
				idx: -3,
				records: []models.Record{
					{ID: -1},
					{ID: -2},
					{ID: -3},
					{ID: 1},
					{ID: 2},
					{ID: 3},
				},
			},
		},
		{
			name: "no new records",
			fields: fields{
				records: []models.Record{
					{ID: 1},
					{ID: 2},
					{ID: 3},
				},
				cryptHasher: nil,
				freeIdx:     0,
			},
			args: args{
				records: nil,
			},
			want: want{
				err: assert.NoError,
				idx: 0,
				records: []models.Record{
					{ID: 1},
					{ID: 2},
					{ID: 3},
				},
			},
		},
		{
			name: "no new records 2",
			fields: fields{
				records: []models.Record{
					{ID: 1},
					{ID: 2},
					{ID: 3},
				},
				cryptHasher: nil,
				freeIdx:     0,
			},
			args: args{
				records: []models.Record{},
			},
			want: want{
				err: assert.NoError,
				idx: 0,
				records: []models.Record{
					{ID: 1},
					{ID: 2},
					{ID: 3},
				},
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
				freeIdx:     tt.fields.freeIdx,
			}
			tt.want.err(t, s.RestoreRecords(tt.args.records), fmt.Sprintf("RestoreRecords(%v)", tt.args.records))
			if !reflect.DeepEqual(s.records, tt.want.records) {
				t.Errorf("records aren't equal ('%v' != '%v')", s.records, tt.want.records)
			}

			assert.Equal(t, tt.want.idx, s.freeIdx)
		})
	}
}
