package inmemory

import (
	"fmt"
	"testing"

	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/stretchr/testify/assert"
)

func TestStorage_AddRecord(t *testing.T) {
	type fields struct {
		records     []models.Record
		cryptHasher *ska.SKA
		freeIdx     int64
	}
	type args struct {
		record *models.Record
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
				records:     nil,
				cryptHasher: nil,
				freeIdx:     0,
			},
			args: args{
				record: &models.Record{ID: 0, Data: models.Data{RecordType: models.TypeText}},
			},
			want: want{
				records: []models.Record{
					{ID: -1, Data: models.Data{RecordType: models.TypeText}},
				},
				err: assert.NoError,
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
			tt.want.err(t, s.AddRecord(tt.args.record), fmt.Sprintf("AddRecord(%v)", tt.args.record))
			assert.Equal(t, tt.want.records[0].ID, s.records[0].ID)
			assert.Equal(t, tt.want.records[0].Data.RecordType, s.records[0].Data.RecordType)
		})
	}
}
