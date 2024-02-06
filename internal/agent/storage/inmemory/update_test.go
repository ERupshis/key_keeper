package inmemory

import (
	"fmt"
	"testing"
	"time"

	"github.com/erupshis/key_keeper/internal/agent/models"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/stretchr/testify/assert"
)

func TestStorage_UpdateRecord(t *testing.T) {
	type fields struct {
		records     []models.Record
		cryptHasher *ska.SKA
		freeIdx     int64
	}
	type args struct {
		date   time.Time
		record *models.Record
	}
	type want struct {
		recordType models.RecordType
		err        assert.ErrorAssertionFunc
		errData    assert.ComparisonAssertionFunc
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
					{1, models.Data{}, false, time.Date(2023, time.January, 11, 12, 0, 0, 0, time.UTC)},
					{2, models.Data{}, false, time.Now()},
					{3, models.Data{}, false, time.Now()},
				},
				cryptHasher: nil,
				freeIdx:     0,
			},
			args: args{
				date:   time.Date(2023, time.January, 11, 12, 0, 0, 0, time.UTC),
				record: &models.Record{ID: 1, Data: models.Data{RecordType: models.TypeText}},
			},
			want: want{
				recordType: models.TypeText,
				err:        assert.NoError,
				errData:    assert.NotEqual,
			},
		},
		{
			name: "missing id in storage",
			fields: fields{
				records: []models.Record{
					{1, models.Data{}, false, time.Date(2023, time.January, 11, 12, 0, 0, 0, time.UTC)},
					{2, models.Data{}, false, time.Now()},
					{3, models.Data{}, false, time.Now()},
				},
				cryptHasher: nil,
				freeIdx:     0,
			},
			args: args{
				date:   time.Date(2023, time.January, 11, 12, 0, 0, 0, time.UTC),
				record: &models.Record{ID: 4, Data: models.Data{RecordType: models.TypeText}},
			},
			want: want{
				recordType: models.TypeText,
				err:        assert.Error,
				errData:    assert.NotEqual,
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

			err := s.UpdateRecord(tt.args.record)
			tt.want.err(t, err, fmt.Sprintf("UpdateRecord(%v)", tt.args.record))
			if err == nil {
				tt.want.errData(t, tt.args.date, s.records[tt.args.record.ID].UpdatedAt)
				assert.NotEqual(t, tt.want.recordType, s.records[tt.args.record.ID].Data.RecordType)
			}
		})
	}
}
