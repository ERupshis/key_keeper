package inmemory

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_DeleteRecord(t *testing.T) {
	type fields struct {
		records     []models.Record
		cryptHasher *ska.SKA
		freeIdx     int64
	}
	type args struct {
		id int64
	}
	type want struct {
		records []models.Record
		err     assert.ErrorAssertionFunc
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
				records: []models.Record{
					{ID: 1, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 2, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 3, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
				},
				cryptHasher: nil,
				freeIdx:     0,
			},
			args: args{
				id: 1,
			},
			want: want{
				records: []models.Record{
					{ID: 1, Data: models.Data{}, Deleted: true, UpdatedAt: time.Now()},
					{ID: 2, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 3, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
				},
				err: assert.NoError,
			},
		},
		{
			name: "negative id",
			fields: fields{
				records: []models.Record{
					{ID: -1, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 2, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 3, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
				},
				cryptHasher: nil,
				freeIdx:     -1,
			},
			args: args{
				id: -1,
			},
			want: want{
				records: []models.Record{
					{ID: 2, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 3, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
				},
				err: assert.NoError,
			},
		},
		{
			name: "wrong id",
			fields: fields{
				records: []models.Record{
					{ID: -1, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 2, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 3, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
				},
				cryptHasher: nil,
				freeIdx:     -1,
			},
			args: args{
				id: -6,
			},
			want: want{
				records: []models.Record{
					{ID: -1, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 2, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
					{ID: 3, Data: models.Data{}, Deleted: false, UpdatedAt: time.Now()},
				},
				err: assert.Error,
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
			tt.want.err(t, s.DeleteRecord(tt.args.id), fmt.Sprintf("DeleteRecord(%v)", tt.args.id))
			if tt.args.id > 0 {
				rec, err := s.GetRecord(tt.args.id)
				require.NoError(t, err)
				assert.True(t, rec.Deleted)
			} else {
				assert.True(t, reflect.DeepEqual(tt.want.records, s.records))
			}
		})
	}
}
